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
	log = logger.New("main")
	mutex = &sync.Mutex{}
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

	gryd := grid.NewGrid(8, dimY+2, true, mask)

	gryd.SetColorAll(color.Black)
	gryd.SetHandler(func(gr *grid.Grid, x, y int, event grid.EventType) {
		if x > 0 && x < dimX + 1 && y > 0 && y < dimY + 1 {
			if event == grid.Pressed {
				flipColor(gryd, x, y, pres.IsPresent())
				updateDisplay(gryd)
			}
		}
	})

	for i := 0; i < dimX+1; i++ {
		gryd.SetColor(i, 0, color.White)
		gryd.SetColor(i, dimY+1, color.White)
	}
	for i := 0; i < dimY+2; i++ {
		gryd.SetColor(0, i, color.White)
		gryd.SetColor(dimX+1, i, color.White)
	}
	setPolarityIndicator(gryd, false)

	go func() {
		changed := make(chan bool, 4)
		pres.Subscribe(changed)
		defer pres.Unsubscribe(changed)
		for range changed {
			invert := pres.IsPresent()
			setPolarityIndicator(gryd, invert)
			for x := 1; x < dimX + 1; x++ {
				for y := 1; y < dimY + 1; y ++ {
					flipColor(gryd, x, y, invert)
				}
			}
			updateDisplay(gryd)
		}
	}()

	updateDisplay(gryd)

	gryd.Connect(pad)
	gryd.Activate()
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

func flipColor(gryd *grid.Grid, x,y int, polarity bool) {
	col := color.FromColor(gryd.Color(x,y))
	onCol, offCol := color.Red, color.Black

	if polarity {
		onCol, offCol = offCol, onCol
	}

	if col.Equal(offCol) {
		gryd.SetColor(x,y, onCol)
	} else {
		gryd.SetColor(x,y, offCol)
	}
}

func updateDisplay(gryd *grid.Grid) {
	msg := display.Message{
		Origin: display.TopLeft,
	}

	for y := 1; y < dimY+1; y++ {
		for x := 1; x < dimX+1; x++ {
			col := color.FromColor(gryd.Color(x,y))
			if col.Equal(color.Black) {
				msg.Pixels += "0"
			} else {
				msg.Pixels += "1"
			}
		}
	}

	log.Info("sending draw command to display")
	_ = disp.Draw(msg)
}

func findDisplay() {
	log := logger.New("scan")

	scanner := discovery.NewScanner(discovery.Display, func(entry discovery.Entry) {
		log.Info("found one display at ", entry.Hostname)
		disp = display.NewRemote(entry)
	}, nil)

	log.Infof("starting scan for displays")
	scanner.Scan()
	<- time.After(time.Millisecond * 500)
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
	<- time.After(time.Millisecond * 500)
	scanner.Close()
}



