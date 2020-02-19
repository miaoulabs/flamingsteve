package main

import (
	"flamingsteve/cmd"
	"image"
	"time"

	"flamingsteve/pkg/ak9753"
	ak9753_presence "flamingsteve/pkg/ak9753/presence"
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
	cmd.SetupLoggers()
	log := logger.New("main")

	muthur.Connect("sensui")

	u := &gui{
		log:     logger.New("gui"),
		changed: make(chan bool),
	}

	wnd := nucular.NewMasterWindowSize(nucular.WindowClosable|nucular.WindowNoScrollbar, "sensor", image.Point{SensorWidth + PropWidth, SensorWidth}, u.renderUi)
	wnd.SetStyle(nstyle.FromTheme(nstyle.DarkTheme, 1.0))
	u.wnd = wnd

	sensorScannner := discovery.NewScanner(
		discovery.Sensor,
		onNewSensor(log, u, wnd),
		onDeleteSensor(log, wnd),
	)
	sensorScannner.Scan()

	detectorScanner := discovery.NewScanner(
		discovery.Detector,
		onNewDetector(log, u, wnd),
		nil,
	)
	detectorScanner.Scan()

	tick := time.NewTicker(time.Second * 60)
	defer tick.Stop()

	// start a scan every 5 second
	go func() {
		for range tick.C {
			sensorScannner.Scan()
			detectorScanner.Scan()
		}
	}()

	go u.updateSensorData()

	wnd.Main()
}

func onDeleteSensor(log *logger.SugaredLogger, wnd nucular.MasterWindow) func(entry discovery.Entry) {
	return func(entry discovery.Entry) {
		log.Infof("removing sensor %s", entry.Id)
		if sensor := sensors.Get(entry.Id); sensor != nil {
			sensor.Device.Close()
			sensors.Delete(entry.Id)
		}
		wnd.Changed()
	}
}

func onNewDetector(log *logger.SugaredLogger, u *gui, wnd nucular.MasterWindow) func(entry discovery.Entry) {
	return func(entry discovery.Entry) {

	}
}

func onNewSensor(log *logger.SugaredLogger, u *gui, wnd nucular.MasterWindow) func(entry discovery.Entry) {
	return func(entry discovery.Entry) {
		if sensor := sensors.Get(entry.Id); sensor == nil {
			log.Infof("new sensor found: %s", entry.Id)

			sub, err := remote.NewSuscriber(entry)
			log.LogIfErr(err)
			if sub == nil {
				return
			}

			//remoteDetector, err := presence.NewSubscriber(entry)
			//log.LogIfErr(err)
			//if remoteDetector == nil {
			//	return
			//}
			//
			//cfgData := remoteDetector.Configs()
			//opts := ak9753_presence.UnmarshalOptions(cfgData)

			detector, _ := ak9753_presence.New(sub, nil)

			dev, err := ak9753.NewMean(sub, detector.Options().Smoothing)
			log.LogIfErr(err)

			sensors.Set(entry.Id, Sensor{
				Ident:          entry,
				Device:         dev,
				LocalDetector:  detector,
			})

			if u.selectedSensorIndex == 0 {
				u.selectedSensorIndex = 1
				u.selectSensor(entry.Id)
			}
		}
		wnd.Changed()
	}
}
