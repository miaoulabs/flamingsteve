package amg8833

import "flamingsteve/pkg/notification"

type Device interface {
	Close()

	Thermistor() float32

	Temperature(x, y int) float32
	Temperatures() [PIXEL_COUNT]float32

	State() State

	/*
	 A true will be pushed every time the sensor's state change
	*/
	notification.Notifier
}
