package amg8833

import "math"

type State struct {
	Thermistor float32
	Pixels     [PIXEL_COUNT]float32
}

func XYtoIndex(x, y int) int {
	if x < 0 || x >= ROW_COUNT || y < 0 || y >= ROW_COUNT {
		return 0
	}
	return x + y*ROW_COUNT
}

func (s *State) Pixel(x, y int) float32 {
	if x < 0 || x >= ROW_COUNT || y < 0 || y >= ROW_COUNT {
		return 0
	}
	return s.Pixels[XYtoIndex(x, y)]
}

func (s *State) Equal(other State) bool {
	for i, temp := range s.Pixels {
		if !sameF32(temp, other.Pixels[i]) {
			return false
		}
	}
	return false
}

func sameF32(f1, f2 float32) bool {
	tolerance := 0.00000001
	return math.Abs(float64(f1)-float64(f2)) < tolerance
}
