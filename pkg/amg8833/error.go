package amg8833

//go:generate go-enum -f=$GOFILE --names

/*
Error x ENUM(
	None = 0
	Param = -1
	Comm = -2
	Other = -128
)
*/
type Error int
