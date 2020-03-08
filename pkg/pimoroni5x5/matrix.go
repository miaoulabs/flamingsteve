package pimoroni5x5

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"sync"
	"time"

	"periph.io/x/periph/conn/i2c"
)

type Display struct {
	dev   i2c.Dev
	mutex sync.Mutex

	pixels []color.Color

	frame uint8
}

func New(bus i2c.Bus, addr i2c.Addr) (*Display, error) {
	d := &Display{
		dev:    i2c.Dev{Addr: uint16(addr), Bus: bus},
		pixels: make([]color.Color, WIDTH*HEIGHT),
	}

	var err error

	if err := d.writeRegister8(CONFIG_BANK, SHUTDOWN_REGISTER, 0x00); err != nil {
		return nil, err
	}

	if err = d.selectFrame(0); err != nil {
		return nil, err
	}

	// turn off all LEDs in the LED control register
	for i := 0; i < 0x11; i++ {
		if err = d.sendUint8(Cmd(i), 0); err != nil {
			return nil, err
		}
	}

	// turn off all LEDs in the blink control register (not really needed)
	for i := 0x12; i < 0x23; i++ {
		if err = d.sendUint8(Cmd(i), 0); err != nil {
			return nil, err
		}
	}

	time.Sleep(time.Millisecond * 10)

	// disable software shutdown
	if err = d.writeRegister8(CONFIG_BANK, SHUTDOWN_REGISTER, 0x01); err != nil {
		return nil, err
	}

	if err = d.AudioSync(false); err != nil {
		return nil, err
	}

	if err = d.writeRegister8(CONFIG_BANK, MODE_REGISTER, REG_CONFIG_PICTUREMODE); err != nil {
		return nil, err
	}

	_ = d.Clear(color.RGBA{})

	return d, nil
}

func (d *Display) Clear(col color.Color) error {
	for i := range d.pixels {
		d.pixels[i] = col
	}
	return d.Show()
}

func (d *Display) DisplayFrame(idx int) error {
	if idx >= MAX_FRAMES {
		d.frame = 0
	}
	return d.writeRegister8(CONFIG_BANK, FRAME_REGISTER, d.frame)
}

func (d *Display) selectFrame(idx uint8) error {
	if idx >= MAX_FRAMES {
		return fmt.Errorf("invalid frame: %v", idx)
	}
	return d.selectBank(Cmd(idx))
}

func (d *Display) Show() error {
	next := (d.frame + 1) % MAX_FRAMES
	var err error

	if err = d.selectFrame(next); err != nil {
		return err
	}

	if err = d.sendUints(Cmd(0x0), ENABLE_PATTERN); err != nil {
		return err
	}

	data := make([]uint8, 144) // controller is 9x16 pixel * 3 colors

	for idx, col := range d.pixels {
		r, g, b, _ := col.RGBA()
		coord := LOOKUP[idx]
		data[coord.r] = uint8(r >> 8)
		data[coord.g] = uint8(g >> 8)
		data[coord.b] = uint8(b >> 8)
	}

	//fmt.Printf("buffer length: %v\n", len(data))

	offset := 0
	for len(data) > 0 {
		lbuf := 32
		if lbuf > len(data) {
			lbuf = len(data)
		}
		buf := data[:lbuf]
		//fmt.Printf("frame: %v, addr: %v, length: %v, data: %v\n", next, offset, len(buf), buf)

		if err = d.sendUints(Cmd(COLOR_OFFSET+offset), buf); err != nil {
			return err
		}
		data = data[lbuf:]
		offset += lbuf
	}

	d.frame = next

	return d.DisplayFrame(int(next))
}

func (d *Display) SetPixel(x, y int, col color.Color) error {
	if x < 0 || x >= d.Width() || y < 0 || y >= d.Height() {
		return fmt.Errorf("invalid pixel coordinate: (%v,%v), dimension: (%v,%v)", x, y, d.Width(), d.Height())
	}
	if y%2 == 1 {
		x = d.Width() - 1 - x
	}

	d.pixels[d.pixelAddr(x, y)] = col
	return nil
}

func (d *Display) pixelAddr(x, y int) int {
	return x + y*d.Width()
}

func (d *Display) SetLEDPWM(bank Cmd, lednum int, pwm uint8) error {
	if lednum >= WIDTH*HEIGHT*3 {
		return fmt.Errorf("invalid pixel number: %v", lednum)
	}
	return d.writeRegister8(bank, Cmd(COLOR_OFFSET+lednum), pwm)
}

func (d *Display) AudioSync(enable bool) error {
	data := uint8(0)
	if enable {
		data = uint8(1)
	}
	return d.writeRegister8(CONFIG_BANK, AUDIOSYNC_REGISTER, data)
}

func (d *Display) Width() int {
	return WIDTH
}

func (d *Display) Height() int {
	return HEIGHT
}

func (d *Display) writeRegister8(bank Cmd, cmd Cmd, data uint8) error {
	err := d.selectBank(bank)
	if err != nil {
		return err
	}
	return d.sendUint8(cmd, data)
}

func (d *Display) selectBank(bank Cmd) error {
	return d.sendUint8(BANK_ADDRESS, byte(bank))
}

func (d *Display) sendUint16(cmd Cmd, data uint16) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := []byte{}
	binary.LittleEndian.PutUint16(bytes, data)
	return d.dev.Tx(append(cmd.toBytes(), bytes...), nil)
}

func (d *Display) sendUint8(cmd Cmd, data uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.dev.Tx(append(cmd.toBytes(), data), nil)
}

func (d *Display) sendUints(cmd Cmd, data []uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.dev.Tx(append(cmd.toBytes(), data...), nil)
}

func (d *Display) readUint8(cmd Cmd) (uint8, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := make([]byte, 1)
	err := d.dev.Tx(cmd.toBytes(), bytes)
	return bytes[0], err
}

func (d *Display) readUint16(cmd Cmd) (uint16, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bytes := make([]byte, 2)
	err := d.dev.Tx(cmd.toBytes(), bytes)
	return binary.LittleEndian.Uint16(bytes), err
}
