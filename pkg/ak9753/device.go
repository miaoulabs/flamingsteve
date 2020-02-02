package ak9753

type Device interface {
	Close()
	DeviceId() (uint8, error)
	CompagnyCode() (uint8, error)
	IR1() (float32, error)
	IR2() (float32, error)
	IR3() (float32, error)
	IR4() (float32, error)
	Temperature() (float32, error)
	All() State

	Subscribe(channel chan<- bool)
}
