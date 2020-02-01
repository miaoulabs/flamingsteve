package ak9753

import "math"

type State struct {
	Ir1, Ir2, Ir3, Ir4 float32
	Temperature        float32
	DeviceId           uint8
	CompagnyCode       uint8
}

func (s State) Equal(other State) bool {
	return sameFloat(s.Temperature, other.Temperature) &&
		sameFloat(s.Ir1, other.Ir1) &&
		sameFloat(s.Ir2, other.Ir2) &&
		sameFloat(s.Ir3, other.Ir3) &&
		sameFloat(s.Ir4, other.Ir4) &&
		s.DeviceId == s.DeviceId &&
		s.CompagnyCode == s.CompagnyCode
}

func sameFloat(f1, f2 float32) bool {
	tolerance := 0.00000001
	return math.Abs(float64(f1) - float64(f2)) < tolerance
}
