package pdetect

//Movement
const (
	FieldCount = 4

	MovementNone     = uint8(0b0000)
	MovementFrom1to3 = 0b0001
	MovementFrom3to1 = 0b0010
	MovementFrom2to4 = 0b0100
	MovementFrom4to2 = 0b1000

	smoothingCount = 6
)


