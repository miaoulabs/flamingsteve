package main

import (
	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/ak9753/remote"
	pdetect "flamingsteve/pkg/presence_detector"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"image"
	"time"
)

var (
	dev      ak9753.Device
	detector *pdetect.Detector
)

const (
	SensorWidth = 800
	PropWidth   = 300
)

func main() {
	var err error
	dev, err = remote.NewSuscriber("protopi")
	if err != nil {
		panic("could not connect to sensor")
	}

	dev, _ = ak9753.NewMean(dev, 5)

	defer dev.Close()

	detector = pdetect.New(dev, nil)
	defer detector.Close()

	u := &ui{
		smoothing: float64(detector.Options().Smoothing),
		presence:  int(detector.Options().PresenceThreshold),
		base:      -200,
		ir:        [ak9753.FieldCount][]float64{},
	}
	wnd := nucular.NewMasterWindowSize(nucular.WindowMovable|nucular.WindowClosable|nucular.WindowNoScrollbar, "sensor", image.Point{SensorWidth + PropWidth, SensorWidth}, u.renderUi)
	wnd.SetStyle(nstyle.FromTheme(nstyle.DarkTheme, 1.0))

	go func() {
		for range detector.PresenceChanged() {
			wnd.Changed()
		}
	}()

	go updateSensorData(wnd, u)

	wnd.Main()
}

func updateSensorData(wnd nucular.MasterWindow, u *ui) {
	maxValue := 60

	changed := make(chan bool)
	dev.Subscribe(changed)

	for range changed {
		if wnd.Closed() { //quit
			return
		}
		wnd.Changed() // force redraw

		u.Lock()
		for i := 0; i < len(u.ir); i++ {
			if len(u.ir[i]) > maxValue {
				u.ir[i] = u.ir[i][1:]
				u.irTime[i] = u.irTime[i][1:]
			}
			u.ir[i] = append(u.ir[i], float64(detector.IR(i)))
			u.irTime[i] = append(u.irTime[i], time.Now())

			//if len(u.ir[i]) > maxValue {
			//	u.ir[i] = u.ir[i][:maxValue]
			//	u.irTime[i] = u.irTime[i][:maxValue]
			//}
			//u.ir[i] = append([]float64{float64(detector.IR(i))}, u.ir[i]...)
			//u.irTime[i] = append([]time.Time{time.Now()}, u.irTime[i]...)
		}
		u.Unlock()
	}
}
