package amg8833

//go:generate go-enum -f=$GOFILE --names --noprefix --prefix=Power

/*
Mode x ENUM(
	Normal = 0
	Sleep = 16
	Standby60s = 32
	Standby10s = 33
)
*/
type PowerMode uint8

//Normal = 0x00
//Sleep = 0x10
//Standby60s = 0x20
//Standby10s = 0x21
