package presence

type smoother struct {
	lastMarked float32
	avgWeigth  float32
	avg        float32
	last       float32
}

func (s *smoother) add(data float32) {
	s.last = data
	s.avg = s.avgWeigth*data + (1-s.avgWeigth)*s.avg
}

func (s *smoother) average() float32 {
	return s.avg
}

func (s *smoother) lastValue() float32 {
	return s.last
}

func (s *smoother) derivative() float32 {
	d := s.avg - s.lastMarked
	s.lastMarked = s.avg
	return d
}
