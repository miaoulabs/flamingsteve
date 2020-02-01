package ak9753

import (
	"encoding/binary"
	"fmt"
	"periph.io/x/periph/conn/i2c"
	"sync"
	"time"
)

type Device struct {
	dev   i2c.Dev
	mutex sync.Mutex
}

func New(bus i2c.Bus, addr i2c.Addr) (*Device, error) {
	d := &Device{
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

func (d *Device) DeviceId() (uint8, error) {
	return d.readUint8(REG_WIA2)
}

func (d *Device) CompagnyCode() (uint8, error) {
	return d.readUint8(REG_WIA1)
}

func (d *Device) DataReady() bool {
	data, err := d.ST1()
	return err == nil && (data&0x01) == 0x01
}

func (d *Device) DataOverRun() bool {
	data, err := d.ST2()
	return err == nil && (data&0x02) == 0x02
}

func (d *Device) ST1() (uint8, error) {
	return d.readUint8(REG_ST1)
}

func (d *Device) ST2() (uint8, error) {
	return d.readUint8(REG_ST2)
}

func (d *Device) INTST() (uint8, error) {
	return d.readUint8(REG_INTST)
}

func (d *Device) RawIR1() (uint16, error) {
	return d.readUint16(REG_IR1L)
}

func (d *Device) RawIR2() (uint16, error) {
	return d.readUint16(REG_IR2L)
}

func (d *Device) RawIR3() (uint16, error) {
	return d.readUint16(REG_IR3L)
}

func (d *Device) RawIR4() (uint16, error) {
	return d.readUint16(REG_IR4L)
}

func (d *Device) IR1() (float32, error) {
	raw, err := d.RawIR1()
	return toFloat(raw), err
}

func (d *Device) IR2() (float32, error) {
	raw, err := d.RawIR2()
	return toFloat(raw), err
}

func (d *Device) IR3() (float32, error) {
	raw, err := d.RawIR3()
	return toFloat(raw), err
}

func (d *Device) IR4() (float32, error) {
	raw, err := d.RawIR4()
	return toFloat(raw), err
}

func (d *Device) RawTemperature() (uint16, error) {
	return d.readUint16(REG_TMPL)
}

func (d *Device) Temperature() (float32, error) {
	raw, err := d.RawTemperature()
	raw >>= 6 // Temp is 10-bit. TMPL0:5 fixed at 0
	temperature := 26.75 + (float32(raw) * 0.125)
	return temperature, err
}

func (d *Device) Temperature_F() (float32, error) {
	temp, err := d.Temperature()
	return temp*1.8 + 32.0, err
}

func (d *Device) ETH13H() (int16, error) {
	return d.readInt16(REG_ETH13H_LSB)
}

func (d *Device) ETH13L() (int16, error) {
	return d.readInt16(REG_ETH13L_LSB)
}

func (d *Device) ETH24H() (int16, error) {
	return d.readInt16(REG_ETH24H_LSB)
}

func (d *Device) ETH24L() (int16, error) {
	return d.readInt16(REG_ETH24L_LSB)
}

func (d *Device) EHYS13() (uint8, error) {
	return d.readUint8(REG_EHYS13)
}

func (d *Device) EHYS24() (uint8, error) {
	return d.readUint8(REG_EHYS24)
}

func (d *Device) EINTE() (uint8, error) {
	return d.readUint8(REG_EINTEN)
}

func (d *Device) ECNTL1() (uint8, error) {
	return d.readUint8(REG_ECNTL1)
}

func (d *Device) CNTL2() (uint8, error) {
	return d.readUint8(REG_CNTL2)
}

func (d *Device) ETHpAtoRaw(pA float32) int16 {
	raw := (int16)(pA / 3.4877)
	if raw > 2047 {
		raw = 2047
	}
	if raw < -2048 {
		raw = -2048
	}
	return raw
}

func (d *Device) SetETH13L(value int16) error {
	return d.sendInt16(REG_ETH13L_LSB, value)
}

func (d *Device) SetETH24H(value int16) error {
	return d.sendInt16(REG_ETH24H_LSB, value)
}

func (d *Device) SetETH24L(value int16) error {
	return d.sendInt16(REG_ETH24L_LSB, value)
}

func (d *Device) EHYSpAtoRaw(pA float32) uint8 {
	raw := (uint16)(pA / 3.4877)
	if raw > 31 {
		raw = 31
	}
	return uint8(raw)
}

func (d *Device) SetEHYS13(val uint8) error {
	return d.sendUint8(REG_EHYS13, val)
}

func (d *Device) SetEHYS24(val uint8) error {
	return d.sendUint8(REG_EHYS24, val)
}

func (d *Device) SetEINTEN(val uint8) error {
	return d.sendUint8(REG_EINTEN, val)
}

func (d *Device) SetECNTL1(val uint8) error {
	return d.sendUint8(REG_ECNTL1, val)
}

func (d *Device) SoftReset() error {
	return d.sendUint8(REG_CNTL2, 0xFF)
}

func (d *Device) StartNextSample() error {
	_, err := d.ST2()
	return err
}

func (d *Device) readUint8(cmd Cmd) (uint8, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := make([]byte, 1)
	err := d.dev.Tx(cmd.toBytes(), bytes)
	return bytes[0], err
}

func (d *Device) readUint16(cmd Cmd) (uint16, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := make([]byte, 2)
	err := d.dev.Tx(cmd.toBytes(), bytes)
	return binary.LittleEndian.Uint16(bytes), err
}

func (d *Device) readInt16(cmd Cmd) (int16, error) {
	raw, err := d.readUint16(cmd)
	return int16(raw), err
}

func (d *Device) sendUint8(cmd Cmd, data uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.dev.Tx(append(cmd.toBytes(), data), nil)
}

func (d *Device) sendInt16(cmd Cmd, data int16) error {
	return d.sendUint16(cmd, uint16(data))
}

func (d *Device) sendUint16(cmd Cmd, data uint16) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := []byte{}
	binary.LittleEndian.PutUint16(bytes, data)
	return d.dev.Tx(append(cmd.toBytes(), bytes...), nil)
}

func toFloat(iVal uint16) float32 {
	//return 14286.8 * float32(iVal) / 32768.0
	return (32767 / 2) * float32(iVal) / 32767
}
