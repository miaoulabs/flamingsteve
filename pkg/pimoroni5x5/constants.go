package pimoroni5x5

const (
	I2C_DEFAULT_ADDRESS   = 0x74
	I2C_ALTERNATE_ADDRESS = 0x77

	MODE_REGISTER      = 0x00
	FRAME_REGISTER     = 0x01
	AUTOPLAY1_REGISTER = 0x02
	AUTOPLAY2_REGISTER = 0x03
	BLINK_REGISTER     = 0x05
	AUDIOSYNC_REGISTER = 0x06
	SHUTDOWN_REGISTER  = 0x0A

	CONFIG_BANK  = 0x0B // helpfully called 'page nine'
	BANK_ADDRESS = 0xFD

	REG_CONFIG_PICTUREMODE   = 0x00
	REG_CONFIG_AUTOPLAYMODE  = 0x08
	REG_CONFIG_AUDIOPLAYMODE = 0x18

	ENABLE_OFFSET = 0x00
	BLINK_OFFSET  = 0x12
	COLOR_OFFSET  = 0x24

	MAX_FRAMES = 8

	WIDTH  = 5
	HEIGHT = 5
)

var (
	ENABLE_PATTERN = []uint8{
		0b00000000, 0b10111111,
		0b00111110, 0b00111110,
		0b00111111, 0b10111110,
		0b00000111, 0b10000110,
		0b00110000, 0b00110000,
		0b00111111, 0b10111110,
		0b00111111, 0b10111110,
		0b01111111, 0b11111110,
		0b01111111, 0b00000000,
	}

	LOOKUP = []lookupCoord{
		{118, 69, 85},
		{117, 68, 101},
		{116, 84, 100},
		{115, 83, 99},
		{114, 82, 98},
		{113, 81, 97},
		{112, 80, 96},
		{134, 21, 37},
		{133, 20, 36},
		{132, 19, 35},
		{131, 18, 34},
		{130, 17, 50},
		{129, 33, 49},
		{128, 32, 48},
		{127, 47, 63},
		{121, 41, 57},
		{122, 25, 58},
		{123, 26, 42},
		{124, 27, 43},
		{125, 28, 44},
		{126, 29, 45},
		{15, 95, 111},
		{8, 89, 105},
		{9, 90, 106},
		{10, 91, 107},
		{11, 92, 108},
		{12, 76, 109},
		{13, 77, 93},
	}
)

type lookupCoord struct {
	r, g, b int
}

type Cmd byte

func (c Cmd) toBytes() []byte {
	return []byte{byte(c)}
}
