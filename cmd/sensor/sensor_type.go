package main

//go:generate go-enum -f=$GOFILE --names --flag

/*
Sensor x ENUM(
	None = -1
	ak9753
	amg8833
)
*/
type SensorType int
