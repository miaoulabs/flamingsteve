package main

type Frame struct {
	id     int
	pixels []bool
}

var frameId = 0

func NewFrame() Frame {
	frameId++
	return Frame{
		pixels: make([]bool, dimX*dimY),
		id:     frameId,
	}
}

func (f *Frame) Copy() Frame {
	cpy := NewFrame()
	copy(cpy.pixels, f.pixels)
	return cpy
}

func (f *Frame) SetPixel(x, y int, lit bool) {
	if x < 0 || x >= dimX || y < 0 || y >= dimY {
		return
	}
	f.pixels[f.idx(x, y)] = lit
}

func (f *Frame) FlipPixel(x, y int) {
	f.SetPixel(x, y, !f.Pixel(x, y))
}

func (f *Frame) Pixel(x, y int) bool {
	if x < 0 || x >= dimX || y < 0 || y >= dimY {
		return false
	}
	return f.pixels[f.idx(x, y)]
}

func (f *Frame) idx(x, y int) int {
	return x + y*dimX
}
