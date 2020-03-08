package main

type Sequence struct {
	index  int
	frames []Frame
}

func NewSequence() Sequence {
	return Sequence{
		frames: []Frame{NewFrame()},
	}
}

func (s *Sequence) Copy() Sequence {
	cpy := Sequence{}
	copy(cpy.frames, s.frames)
	return cpy
}

func (s *Sequence) InsertFrame(before bool) {
	var idx = s.index
	if before {
		s.index++ // current need to stay on the same frame
	} else {
		idx++
	}
	s.frames = append(s.frames[:idx], append([]Frame{NewFrame()}, s.frames[idx:]...)...)
	log.Infof("inserting frame at %v, new length: %d", idx, len(s.frames))
}

func (s *Sequence) Next() {
	if s.index < len(s.frames)-1 {
		s.index++
	}
}

func (s *Sequence) Previous() {
	if s.index > 0 {
		s.index--
	}
}

func (s *Sequence) Current() *Frame {
	return &s.frames[s.index]
}

func (s *Sequence) Play() {

}

func (s *Sequence) Stop() {

}

func (s *Sequence) PlayPause() {

}
