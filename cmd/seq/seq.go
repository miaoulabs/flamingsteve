package main

import (
	"flamingsteve/cmd"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/display"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/presence"
	"github.com/draeron/golaunchpad/pkg/device"
	"github.com/draeron/golaunchpad/pkg/grid"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
	"sync"
	"time"
)

const (
	dimX = 5
	dimY = 10
)

var (
	disp *display.Remote
	pres presence.Detector
	log  = logger.New("main")
	ctrl *grid.Grid

	sequence Sequence = NewSequence()

	bank  = [8]*Sequence{}
	mutex = sync.RWMutex{}
)

func main() {
	device.SetLogger(logger.New("device"))
	minimk3.SetLogger(logger.New("minimk3"))
	cmd.SetupLoggers()

	log.Infof("seq starting")
	defer log.Infof("seq stopped")

	muthur.Connect("seq")
	defer muthur.Close()

	findDisplay()
	if disp == nil {
		log.Fatal("could not find any display on the network")
	}

	findDetector()
	if pres == nil {
		log.Fatal("could not find any detector on the network")
	}

	pad, err := minimk3.Open(minimk3.ProgrammerMode)
	cmd.Must(err)
	cmd.Must(pad.Diag())

	mask := launchpad.Mask{
		button.User: true,
	}
	mask = mask.MergePreset(
		launchpad.MaskRows,
		launchpad.MaskArrows,
	)

	ctrl = grid.NewGrid(8, dimY+2, false, mask)

	ctrl.SetColorAll(color.Black)
	ctrl.SetHandler(handleGrid)

	ctrl.Layout.SetHandlerHold(launchpad.RowHold, handleRowHold)
	ctrl.Layout.SetHandler(launchpad.RowReleased, handleRowReleased)
	ctrl.Layout.SetHandlerHold(launchpad.ArrowHold, handleArrowHold)
	ctrl.Layout.SetHoldTimer(launchpad.ArrowHold, time.Millisecond*100)
	ctrl.Layout.SetHandler(launchpad.ArrowReleased, handleArrowReleased)

	// draw contour
	for i := 0; i < dimX+1; i++ {
		ctrl.SetColor(i, 0, color.White)
		ctrl.SetColor(i, dimY+1, color.White)
	}
	for i := 0; i < dimY+2; i++ {
		ctrl.SetColor(0, i, color.White)
		ctrl.SetColor(dimX+1, i, color.White)
	}

	go func() {
		changed := make(chan bool, 4)
		pres.Subscribe(changed)
		defer pres.Unsubscribe(changed)
		for range changed {
			invert := pres.IsPresent()
			setPolarityIndicator(ctrl, invert)

			mutex.RLock()
			update()
			mutex.RUnlock()
		}
	}()

	setPolarityIndicator(ctrl, false)
	update()

	ctrl.Connect(pad)
	ctrl.Activate()
	cmd.WaitForCtrlC()
}

func setPolarityIndicator(gryd *grid.Grid, polarity bool) {
	if polarity {
		gryd.SetColor(0, dimY/2, color.Orange)
		gryd.SetColor(dimX+1, dimY/2, color.Orange)
	} else {
		gryd.SetColor(0, dimY/2, color.Cyan)
		gryd.SetColor(dimX+1, dimY/2, color.Cyan)
	}
}

func update() {
	updateGrid()
	updateDisplay()
}

func updateGrid() {
	onCol, offCol := color.Red, color.Black
	if pres.IsPresent() {
		onCol, offCol = offCol, onCol
	}

	frame := sequence.Current()

	for x := 0; x < dimX; x++ {
		for y := 0; y < dimY; y++ {
			if frame.Pixel(x, y) {
				ctrl.SetColor(x+1, y+1, onCol)
			} else {
				ctrl.SetColor(x+1, y+1, offCol)
			}
		}
	}

	if ctrl.CanPanUp() {
		ctrl.Layout.SetColor(button.Up, color.White)
	} else {
		ctrl.Layout.SetColor(button.Up, color.Red)
	}

	if ctrl.CanPanDown() {
		ctrl.Layout.SetColor(button.Down, color.White)
	} else {
		ctrl.Layout.SetColor(button.Down, color.Red)
	}

	if sequence.index > 0 {
		ctrl.Layout.SetColor(button.Left, color.White)
	} else {
		ctrl.Layout.SetColor(button.Left, color.Red)
	}

	if sequence.index < len(sequence.frames)-1 {
		ctrl.Layout.SetColor(button.Right, color.White)
	} else {
		ctrl.Layout.SetColor(button.Right, color.Red)
	}

	for i := 0; i < dimX+2; i++ {
		col := color.White
		if i == sequence.index%(dimX+2) {
			col = color.Blue
		}
		ctrl.SetColor(i, 0, col)
		ctrl.SetColor(i, dimY+1, col)
	}

	for idx, seq := range bank {
		if seq == nil {
			ctrl.Layout.SetColor(button.Row1+button.Button(idx), color.Yellow)
		} else {
			ctrl.Layout.SetColor(button.Row1+button.Button(idx), color.Orange)
		}
	}

	ctrl.Layout.SetColor(button.StopSoloMute, color.YellowGreen)
}

func updateDisplay() {
	msg := display.Message{
		Origin: display.TopLeft,
	}

	on, off := "1", "0"
	if pres.IsPresent() {
		on, off = off, on
	}

	frame := sequence.Current()

	for y := 0; y < dimY; y++ {
		for x := 0; x < dimX; x++ {
			if frame.Pixel(x, y) {
				msg.Pixels += on
			} else {
				msg.Pixels += off
			}
		}
	}

	//log.Info("sending draw command to display")
	_ = disp.Draw(msg)
}

func findDisplay() {
	log := logger.New("scan")

	scanner := discovery.NewScanner(discovery.Display, func(entry discovery.Entry) {
		log.Info("found one display at ", entry.Hostname)
		var err error
		disp, err = display.NewRemote(entry)
		if err != nil {
			log.Errorf("fail to create remote display: %v", err)
		}
	}, nil)

	log.Infof("starting scan for displays")
	scanner.Scan()
	<-time.After(time.Millisecond * 500)
	scanner.Close()
}

func findDetector() {
	log := logger.New("scan")

	scanner := discovery.NewScanner(discovery.Detector, func(entry discovery.Entry) {
		log.Info("found one detector at ", entry.Hostname)
		var err error
		pres, err = presence.NewSubscriber(entry)
		cmd.Must(err)
	}, nil)

	log.Infof("starting scan for displays")
	scanner.Scan()
	<-time.After(time.Millisecond * 500)
	scanner.Close()
}
