package amg8833

//go:generate go-enum -f=$GOFILE --names --noprefix --prefix=Interrupt

/*
InterruptMode x ENUM(
	Disabled       = 0
	Active         = 1
	AbsoluteValue  = 2
)
*/
type InterruptMode uint8
