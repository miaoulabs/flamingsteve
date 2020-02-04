package main

import (
	"image"
	"time"

	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/ak9753/presence"
	"flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/muthur"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/draeron/gopkgs/logger"
)

var (
	dev      ak9753.Device
	detector *presence.Detector
)

const (
	SensorWidth = 800
	PropWidth   = 300
)

func main() {
	log := logger.New("main")
	presence.SetLogger(logger.New("detector"))
	remote.SetLogger(logger.New("remote"))
	ak9753.SetLogger(logger.New("ak9753"))

	var mothers muthur.Servers
	for mothers == nil || len(mothers) == 0 {
		mothers = muthur.ResolveServers(time.Second * 2)
	}

	// we only need one
	log.Infof("connecting to muthur on host '%s'", mothers[0].HostName)

	var err error
	dev, err = remote.NewSuscriber(mothers[0].HostName, "sensor-ui")
	if err != nil {
		panic("could not connect to sensor")
	}

	dev, _ = ak9753.NewMean(dev, 5)

	defer dev.Close()

	detector = presence.New(dev, nil)
	defer detector.Close()

	u := &ui{
		smoothing: float64(detector.Options().Smoothing),
		presence:  int(detector.Options().PresenceThreshold),
		base:      -200,
		ir:        [ak9753.FieldCount][]float64{},
		log:       logger.New("ui"),
	}

	wnd := nucular.NewMasterWindowSize(nucular.WindowMovable|nucular.WindowClosable|nucular.WindowNoScrollbar, "sensor", image.Point{SensorWidth + PropWidth, SensorWidth}, u.renderUi)
	wnd.SetStyle(nstyle.FromTheme(nstyle.DarkTheme, 1.0))

	go func() {
		for range detector.PresenceChanged() {
			wnd.Changed()
		}
	}()

	go u.updateSensorData(wnd)

	wnd.Main()
}
