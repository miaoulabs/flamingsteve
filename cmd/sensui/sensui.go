package main

import (
	"image"
	"time"

	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/ak9753/presence"
	"flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/draeron/gopkgs/logger"
)

var (
	sensors SensorsMap
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
	muthur.SetLogger(logger.New("muthur"))

	muthur.Connect("sensui")

	u := &gui{
		//smoothing:        float64(detector.Options().Smoothing),
		//presence:         int(detector.Options().PresenceThreshold),
		log:                 logger.New("gui"),
		selectedSensorIndex: 0,
		changed:             make(chan bool),
	}

	wnd := nucular.NewMasterWindowSize(nucular.WindowClosable|nucular.WindowNoScrollbar, "sensor", image.Point{SensorWidth + PropWidth, SensorWidth}, u.renderUi)
	wnd.SetStyle(nstyle.FromTheme(nstyle.DarkTheme, 1.0))
	u.wnd = wnd

	//go func() {
	//	for range detector.PresenceChanged() {
	//		wnd.Changed()
	//	}
	//}()

	scanner := discovery.NewScanner(
		discovery.Sensor,
		logger.New("scanner"),
		func(entry discovery.Entry) {
			if sensor := sensors.Get(entry.Id); sensor == nil {
				log.Infof("new sensor found: %s", entry.Id)

				sub, err := remote.NewSuscriber(&entry)
				log.LogIfErr(err)
				if sub == nil {
					return
				}

				detector, _ := presence.New(sub, nil)

				dev, err := ak9753.NewMean(sub, detector.Options().Smoothing)
				log.LogIfErr(err)

				sensors.Set(entry.Id, Sensor{
					Ident:    entry,
					Device:   dev,
					Detector: detector,
				})

				if u.selectedSensorIndex == 0 {
					u.selectedSensorIndex = 1
					u.selectSensor(entry.Id)
				}
			}
			wnd.Changed()
		},
		func(entry discovery.Entry) {
			log.Infof("removing sensor %s", entry.Id)
			if sensor := sensors.Get(entry.Id); sensor != nil {
				sensor.Device.Close()
				sensors.Delete(entry.Id)
			}
			wnd.Changed()
		},
	)
	scanner.Scan()

	// start a scan every 5 second
	go func() {
		tick := time.NewTicker(time.Second * 5)
		defer tick.Stop()
		for range tick.C {
			scanner.Scan()
		}
	}()

	//detector = presence.New(dev, nil)
	//defer detector.Close()

	go u.updateSensorData()

	wnd.Main()
}


