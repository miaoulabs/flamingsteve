package main

import (
	"github.com/draeron/golaunchpad/pkg/grid"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkgs/color"
	"time"
)

// Short press (release) on rows will play / stop a sequence
func handleRowReleased(layout *launchpad.BasicLayout, btn button.Button) {
	mutex.Lock()
	defer mutex.Unlock()

	if layout.IsHold(btn, launchpad.DefaultHoldDuration) {
		update()
		return
	}

	idx := btn - button.Row1

	if bank[idx] == nil || idx == 7 {
		sequence = NewSequence()
	} else {
		sequence = bank[idx].Copy()
		sequence.index = 0
	}
	update()
}

// Holding row will assign the sequence to a row (first tick)
func handleRowHold(layout *launchpad.BasicLayout, btn button.Button, first bool) {
	idx := btn - button.Row1
	if !first || idx == 7 {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	cpy := sequence.Copy()
	bank[idx] = &cpy
	update()
	layout.SetColor(btn, color.Green)
}

// up/down => pan, l/R => change frame left or right
func handleArrowReleased(layout *launchpad.BasicLayout, btn button.Button) {
	mutex.Lock()
	defer mutex.Unlock()

	long := layout.HoldTime(btn) > time.Millisecond*500

	switch btn {
	case button.Up:
		ctrl.PanUp()
	case button.Down:
		ctrl.PanDown()

	case button.Left:
		if long {
			sequence.InsertFrame(true)
		}
		sequence.Previous()
	case button.Right:
		if long {
			sequence.InsertFrame(false)
		}
		sequence.Next()
	}

	log.Infof("current frame: %v", sequence.index)
	update()
}

// Holding up/down => pan, holding L/R will create a new frame (first tick)
func handleArrowHold(layout *launchpad.BasicLayout, btn button.Button, first bool) {
	mutex.Lock()
	defer mutex.Unlock()

	switch btn {
	case button.Up:
		ctrl.PanUp()
	case button.Down:
		ctrl.PanDown()
	}

	update()
	layout.SetColor(btn, color.Green)
}

func handleGrid(gryd *grid.Grid, x, y int, event grid.EventType) {
	mutex.Lock()
	defer mutex.Unlock()

	if x < 1 || x > dimX || y < 1 || y > dimY || event != grid.Pressed {
		return
	}
	sequence.Current().FlipPixel(x-1, y-1)
	update()
}
