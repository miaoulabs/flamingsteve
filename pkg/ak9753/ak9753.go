package ak9753

import (
	"encoding/binary"
	"fmt"
	"periph.io/x/periph/conn/i2c"
	"sync"
	"time"
)

type Physical struct {
	dev   i2c.Dev
	mutex sync.Mutex
}

func New(bus i2c.Bus, addr i2c.Addr) (*Physical, error) {
	d := &Physical{
		dev: i2c.Dev{Addr: uint16(addr), Bus: bus},
	}

	// wait 3 ms
	time.Sleep(time.Millisecond * 3)

	id, err := d.DeviceId()
	if err != nil {
		return nil, err
	}

	err = d.SoftReset()
	if err != nil {
		return nil, err
	}

	if id != DEVICE_ID {
		return nil, fmt.Errorf("invalid device id, expecting 0x%x, found 0x%x", DEVICE_ID, id)
	}

	//set mode and filter freq
	if err := d.SetECNTL1(uint8(AK975X_FREQ_8_8HZ<<3) | AK975X_MODE_0); err != nil {
		return nil, err
	}

	//enable interrupt
	if err = d.SetEINTEN(0x1f); err != nil { //enable all interrupts
		return nil, err
	}

	return d, nil
}

func (d *Physical) DeviceId() (uint8, error) {
	return d.readUint8(REG_WIA2)
}

func (d *Physical) CompagnyCode() (uint8, error) {
	return d.readUint8(REG_WIA1)
}

func (d *Physical) Model() (string, error) {
	model, err := d.readUint8(REG_INFO1)
	if err != nil {
		return "", fmt.Errorf("could not determine mode", model)
	}
	if model == SENSOR_VERSION_AK9750 {
		return "ak9750", nil
	} else if model == SENSOR_VERSION_AK9753 {
		return "ak9753", nil
	} else {
		return "unknown", nil
	}
}

func (d *Physical) DataReady() bool {
	data, err := d.ST1()
	return err == nil && (data & (1 << 0)) != 0 //Bit 0 is DRDY
}

func (d *Physical) DataOverRun() bool {
	data, err := d.ST1()
	return err == nil && (data & 1 << 1) != 0 //Bit 1 is DOR
}

func (d *Physical) ST1() (uint8, error) {
	return d.readUint8(REG_ST1)
}

func (d *Physical) ST2() (uint8, error) {
	return d.readUint8(REG_ST2)
}

func (d *Physical) INTST() (uint8, error) {
	return d.readUint8(REG_INTST)
}

func (d *Physical) RawIR1() (uint16, error) {
	return d.readUint16(REG_IR1L)
}

func (d *Physical) RawIR2() (uint16, error) {
	return d.readUint16(REG_IR2L)
}

func (d *Physical) RawIR3() (uint16, error) {
	return d.readUint16(REG_IR3L)
}

func (d *Physical) RawIR4() (uint16, error) {
	return d.readUint16(REG_IR4L)
}

func (d *Physical) IR1() (float32, error) {
	raw, err := d.RawIR1()
	return toFloat(raw), err
}

func (d *Physical) IR2() (float32, error) {
	raw, err := d.RawIR2()
	return toFloat(raw), err
}

func (d *Physical) IR3() (float32, error) {
	raw, err := d.RawIR3()
	return toFloat(raw), err
}

func (d *Physical) IR4() (float32, error) {
	raw, err := d.RawIR4()
	return toFloat(raw), err
}

func (d *Physical) RawTemperature() (uint16, error) {
	return d.readUint16(REG_TMPL)
}

func (d *Physical) Temperature() (float32, error) {
	raw, err := d.RawTemperature()
	iraw := int16(raw) >> 6 // Temp is 10-bit. TMPL0:5 fixed at 0
	temperature := 26.75 + (float32(iraw) * 0.125)
	return temperature, err
}

func (d *Physical) Temperature_F() (float32, error) {
	temp, err := d.Temperature()
	return temp*1.8 + 32.0, err
}

func (d *Physical) ETH13H() (int16, error) {
	return d.readInt16(REG_ETH13H_LSB)
}

func (d *Physical) ETH13L() (int16, error) {
	return d.readInt16(REG_ETH13L_LSB)
}

func (d *Physical) ETH24H() (int16, error) {
	return d.readInt16(REG_ETH24H_LSB)
}

func (d *Physical) ETH24L() (int16, error) {
	return d.readInt16(REG_ETH24L_LSB)
}

func (d *Physical) EHYS13() (uint8, error) {
	return d.readUint8(REG_EHYS13)
}

func (d *Physical) EHYS24() (uint8, error) {
	return d.readUint8(REG_EHYS24)
}

func (d *Physical) EINTE() (uint8, error) {
	return d.readUint8(REG_EINTEN)
}

func (d *Physical) ECNTL1() (uint8, error) {
	return d.readUint8(REG_ECNTL1)
}

func (d *Physical) CNTL2() (uint8, error) {
	return d.readUint8(REG_CNTL2)
}

func (d *Physical) ETHpAtoRaw(pA float32) int16 {
	raw := (int16)(pA / 3.4877)
	if raw > 2047 {
		raw = 2047
	}
	if raw < -2048 {
		raw = -2048
	}
	return raw
}

func (d *Physical) SetETH13L(value int16) error {
	return d.sendInt16(REG_ETH13L_LSB, value)
}

func (d *Physical) SetETH24H(value int16) error {
	return d.sendInt16(REG_ETH24H_LSB, value)
}

func (d *Physical) SetETH24L(value int16) error {
	return d.sendInt16(REG_ETH24L_LSB, value)
}

func (d *Physical) EHYSpAtoRaw(pA float32) uint8 {
	raw := (uint16)(pA / 3.4877)
	if raw > 31 {
		raw = 31
	}
	return uint8(raw)
}

func (d *Physical) SetEHYS13(val uint8) error {
	return d.sendUint8(REG_EHYS13, val)
}

func (d *Physical) SetEHYS24(val uint8) error {
	return d.sendUint8(REG_EHYS24, val)
}

func (d *Physical) SetEINTEN(val uint8) error {
	return d.sendUint8(REG_EINTEN, val)
}

func (d *Physical) SetECNTL1(val uint8) error {
	return d.sendUint8(REG_ECNTL1, val)
}

func (d *Physical) SoftReset() error {
	return d.sendUint8(REG_CNTL2, 0xFF)
}

func (d *Physical) StartNextSample() error {
	_, err := d.ST2()
	return err
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

func (d *Physical) readInt16(cmd Cmd) (int16, error) {
	raw, err := d.readUint16(cmd)
	return int16(raw), err
}

func (d *Physical) sendUint8(cmd Cmd, data uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.dev.Tx(append(cmd.toBytes(), data), nil)
}

func (d *Physical) sendInt16(cmd Cmd, data int16) error {
	return d.sendUint16(cmd, uint16(data))
}

func (d *Physical) sendUint16(cmd Cmd, data uint16) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := []byte{}
	binary.LittleEndian.PutUint16(bytes, data)
	return d.dev.Tx(append(cmd.toBytes(), bytes...), nil)
}

func toFloat(val uint16) float32 {
	ival := int16(val)
	return float32(ival)
}
