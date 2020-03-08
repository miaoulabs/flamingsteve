package amg8833

import (
	"encoding/binary"
	"sync"
	"time"

	"periph.io/x/periph/conn/i2c"
)

type Physical struct {
	dev   i2c.Dev
	mutex sync.Mutex
}

func New(bus i2c.Bus, addr i2c.Addr) (*Physical, error) {
	p := &Physical{
		dev: i2c.Dev{Addr: uint16(addr), Bus: bus},
	}

	var err error

	err = p.SetPowerMode(PowerNormal)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 30)

	//err = p.SetUpperLimit(DEFAULT_UPPER_LIMIT)
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = p.SetLowerLimit(DEFAULT_LOWER_LIMIT)
	//if err != nil {
	//	return nil, err
	//}

	//err = p.SetHysteresis(DEFAULT_HYSTERESIS)
	//if err != nil {
	//	return nil, err
	//}

	err = p.SetInterruptMode(InterruptDisabled)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 10)

	err = p.ClearStatus(CLEAR_ALL_STATUS)
	if err != nil {
		return nil, err
	}

	err = p.ResetFlags(INIT_RESET_VALUE)
	if err != nil {
		return nil, err
	}

	err = p.SetFPS(FPS_10)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 100)

	return p, nil
}

func (p *Physical) SetFPS(fps uint8) error {
	return p.sendUint8(FRAME_RATE_ADDR, fps)
}

func (p *Physical) SetPowerMode(mode PowerMode) error {
	return p.sendUint8(POWER_CONTROL_REG_ADDR, uint8(mode))
}

func (p *Physical) SetInterruptMode(mode InterruptMode) error {
	return p.sendUint8(INTERRUPT_CONTROL_REG_ADDR, uint8(mode))
}

func (p *Physical) InterruptMode() (bool, error) {
	val, err := p.readUint8(STATUS_REG_ADDR)
	if err != nil {
		return false, err
	}
	return val&0x02 != 0, nil
}

func (p *Physical) ClearStatus(value uint8) error {
	return p.sendUint8(STATUS_CLEAR_REG_ADDR, value)
}

func (p *Physical) PixelsInterruptStatus(status [ROW_COUNT]uint8) error {
	var err error
	for i := 0; i < ROW_COUNT; i++ {
		status[i], err = p.readUint8(Cmd(INTERRUPT_TABLE_1_8_REG_ADDR + i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Physical) SetUpperLimit(value uint16) error {
	return p.sendUint16(INT_LEVEL_REG_ADDR_HL, value)
}

func (p *Physical) SetLowerLimit(value uint16) error {
	return p.sendUint16(INT_LEVEL_REG_ADDR_LL, value)
}

func (p *Physical) SetHysteresis(value uint16) error {
	return p.sendUint16(INT_LEVEL_REG_ADDR_YSL, value)
}

/*
	INIT_RESET_VALUE or FLAG_RESET_VALUE
*/
func (p *Physical) ResetFlags(value uint8) error {
	return p.sendUint8(RESET_REG_ADDR, value)
}

func (p *Physical) PixelTemperatureRaw() ([PIXEL_COUNT]uint16, error) {
	reading := [PIXEL_COUNT]uint16{}
	temps := [PIXEL_COUNT]uint16{}

	var err error
	for i := 0; i < PIXEL_COUNT; i++ {
		reading[i], err = p.readUint16(Cmd(TEMPERATURE_REG_ADDR_L + 2*i))
		if err != nil {
			return temps, err
		}
	}

	// flip on Y axis
	for x := 0; x < ROW_COUNT; x++ {
		for y := 0; y < ROW_COUNT; y++ {
			temps[x+y*ROW_COUNT] = reading[x+(ROW_COUNT-1-y)*ROW_COUNT]
		}
	}

	return temps, nil
}

func (p *Physical) PixelTemperature() ([PIXEL_COUNT]float32, error) {
	temps := [PIXEL_COUNT]float32{}
	raw, err := p.PixelTemperatureRaw()
	if err != nil {
		return temps, err
	}

	for i, rawtemp := range raw {
		temps[i] = int12ToFloat(rawtemp) * PIXEL_TEMP_CONVERSION
	}
	return temps, nil
}

func (p *Physical) ThermistorRaw() (uint16, error) {
	return p.readUint16(THERMISTOR_REG_ADDR_L)
}

func (p *Physical) Thermistor() (float32, error) {
	raw, err := p.ThermistorRaw()
	if err != nil {
		return 0, err
	}
	return signedMag12ToFloat(raw) * THERMISTOR_CONVERSION, nil
}

/*!
  @brief  convert a 12-bit signed magnitude value to a floating point number
  @param  val the 12-bit signed magnitude value to be converted
  @returns the converted floating point value
*/
func signedMag12ToFloat(raw uint16) float32 {
	//take first 11 bits as absolute val
	absVal := uint16(raw & 0x7FF)

	if raw&0x800 != 0 {
		return float32(-absVal)
	} else {
		return float32(absVal)
	}
}

/*!
  @brief  convert a 12-bit integer two's complement value to a floating point number
  @param  val the 12-bit integer  two's complement value to be converted
  @returns the converted floating point value
*/
func int12ToFloat(val uint16) float32 {
	sVal := int16(val << 4)   //shift to left so that sign bit of 12 bit integer number is placed on sign bit of 16 bit signed integer number
	return float32(sVal >> 4) //shift back the signed number, return converts to float
}

func (d *Physical) sendUint16(cmd Cmd, data uint16) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := []byte{}
	binary.LittleEndian.PutUint16(bytes, data)
	return d.dev.Tx(append(cmd.toBytes(), bytes...), nil)
}

func (d *Physical) sendUint8(cmd Cmd, data uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.dev.Tx(append(cmd.toBytes(), data), nil)
}

func (d *Physical) readUint8(cmd Cmd) (uint8, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := make([]byte, 1)
	err := d.dev.Tx(cmd.toBytes(), bytes)
	return bytes[0], err
}

func (d *Physical) readUint16(cmd Cmd) (uint16, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := make([]byte, 2)
	err := d.dev.Tx(cmd.toBytes(), bytes)
	return binary.LittleEndian.Uint16(bytes), err
}
