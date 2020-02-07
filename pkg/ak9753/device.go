package ak9753

import "flamingsteve/pkg/notify"

type Device interface {
	Close()
	DeviceId() (uint8, error)
	CompagnyCode() (uint8, error)

	IR(idx int) (float32, error)
	IR1() (float32, error)
	IR2() (float32, error)
	IR3() (float32, error)
	IR4() (float32, error)
	Temperature() (float32, error)
	All() State

	/*
	 A true will be pushed every time the sensor's state change
	*/
	notify.Nofifier
}
