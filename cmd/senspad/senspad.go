package main

import (
	"time"

	"flamingsteve/cmd"
	"flamingsteve/pkg/amg8833"
	"flamingsteve/pkg/amg8833/remote"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/display"
	"flamingsteve/pkg/muthur"

	"github.com/draeron/golaunchpad/pkg/device"
	"github.com/draeron/golaunchpad/pkg/grid"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
	"github.com/wcharczuk/go-chart"
)

var (
	disp   *display.Remote
	log    = logger.New("main")
	ctrl   *grid.Grid
	sensor amg8833.Device
)

func main() {
	device.SetLogger(logger.New("device"))
	minimk3.SetLogger(logger.New("minimk3"))
	cmd.SetupLoggers()

	log.Infof("senspad starting")
	defer log.Infof("senspad stopped")

	muthur.Connect("senspad")
	defer muthur.Close()

	findSensor()

	pad, err := minimk3.Open(minimk3.ProgrammerMode)
	cmd.Must(err)
	cmd.Must(pad.Diag())
	defer pad.Close()

	ctrl = grid.NewGrid(amg8833.ROW_COUNT, amg8833.ROW_COUNT, false, launchpad.Mask{})
	ctrl.SetColorAll(color.Black)

	go update()

	ctrl.Connect(pad)
	ctrl.Activate()
	cmd.WaitForCtrlC()
}

func update() {
	changed := make(chan bool)
	sensor.Subscribe(changed)
	defer close(changed)

	for range changed {
		state := sensor.State()

		for x := 0; x < amg8833.ROW_COUNT; x++ {
			for y := 0; y < amg8833.ROW_COUNT; y++ {
				temp := state.Pixel(x, y)
				col := chart.Jet(float64(temp), 26, 28)
				ctrl.SetColor(amg8833.ROW_COUNT-1-x, amg8833.ROW_COUNT-1-y, col)
			}
		}
	}
}

func findSensor() {
	log := logger.New("scan")

	scanner := discovery.NewScanner(discovery.Sensor, func(entry discovery.Entry) {
		log.Info("found one detector at ", entry.Hostname)
		if entry.Model != amg8833.ModelName {
			return
		}
		var err error
		sensor, err = remote.NewSuscriber(entry)
		cmd.Must(err)

		sensor, err = amg8833.NewMean(sensor, 5)
		cmd.Must(err)
	}, nil)

	log.Infof("starting scan for displays")
	scanner.Scan()
	<-time.After(time.Millisecond * 500)
	scanner.Close()
}
