package main

import (
	"image"
	"strconv"
	"sync"

	"flamingsteve/pkg/display"
	"github.com/aarzilli/nucular"
	"github.com/draeron/gopkgs/color"
	"github.com/fogleman/gg"
)

type Gui struct {
	MainWindow nucular.MasterWindow
	listener   *display.Listener
	pixW       int
	pixH       int
	pixels     []bool
	mutex      sync.RWMutex
}

func NewGui(w,h int) *Gui {
	g := &Gui{
		pixW: w,
		pixH: h,
		pixels: make([]bool, w*h),
	}
	g.listener = display.NewListener("dispsim", "virtual", g.updateDrawing)

	return g
}

func (g *Gui) WindowSize() image.Point {
	return image.Point{
		X: g.pixW * (PixelEdgeDimension + PixelSpacing),
		Y: g.pixH * (PixelEdgeDimension + PixelSpacing) + PixelSpacing/2,
	}
}

func (g *Gui) updateDrawing(drawMsg *display.Message) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if drawMsg.Origin == "" {
		drawMsg.Origin = display.TopLeft
	}

	switch drawMsg.Origin {
	case display.TopLeft:
		for idx := range g.pixels {
			if idx < len(drawMsg.Pixels) {
				g.pixels[idx] = toBool(drawMsg.Pixels[idx])
			} else if drawMsg.ClearOnMissing {
				g.pixels[idx] = false
			}
		}

	case display.BottomRight:
		for idx := range g.pixels {
			invertedIdx := len(g.pixels) - idx - 1
			if idx < len(drawMsg.Pixels) {
				g.pixels[invertedIdx] = toBool(drawMsg.Pixels[idx])
			} else if drawMsg.ClearOnMissing {
				g.pixels[invertedIdx] = false
			}
		}
	}

	g.MainWindow.Changed()
}

func toBool(r byte) bool {
	if r == '1' {
		return true
	} else {
		return false
	}
}

func (g *Gui) render(mw *nucular.Window) {
	mw.Row(PixelEdgeDimension + PixelSpacing/2).Dynamic(g.pixW)
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	for idx, pix := range g.pixels {
		img := image.NewRGBA(image.Rect(0,0, PixelEdgeDimension + PixelSpacing/2, PixelEdgeDimension + PixelSpacing/2))
		ctx := gg.NewContextForRGBA(img)
		ctx.SetColor(mw.WindowStyle().Background)
		ctx.Clear()
		if pix {
			ctx.SetColor(color.Red)
		} else {
			ctx.SetColor(color.Black)
		}
		ctx.DrawRoundedRectangle(PixelSpacing/4, PixelSpacing/4, PixelEdgeDimension, PixelEdgeDimension, PixelSpacing)
		ctx.Fill()

		ctx.SetColor(color.White)
		ctx.DrawStringAnchored(strconv.Itoa(idx+1), PixelEdgeDimension/2, PixelEdgeDimension/2, 0.5, 0.5)

		mw.Image(img)
	}
}
