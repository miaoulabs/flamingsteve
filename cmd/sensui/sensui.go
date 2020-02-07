package main

import (
	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/ak9753/presence"
	"flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/draeron/gopkgs/logger"
	"image"
	"sort"
	"sync"
	"time"
)

var (
	sensors sync.Map
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

	u := &ui{
		//smoothing:        float64(detector.Options().Smoothing),
		//presence:         int(detector.Options().PresenceThreshold),
		log:              logger.New("ui"),
		selectedSensorId: 0,
		changed:          make(chan bool),
	}

	wnd := nucular.NewMasterWindowSize(nucular.WindowMovable|nucular.WindowClosable|nucular.WindowNoScrollbar, "sensor", image.Point{SensorWidth + PropWidth, SensorWidth}, u.renderUi)
	wnd.SetStyle(nstyle.FromTheme(nstyle.DarkTheme, 1.0))
	u.wnd = wnd

	//go func() {
	//	for range detector.PresenceChanged() {
	//		wnd.Changed()
	//	}
	//}()

	scanner := discovery.NewScanner(
		"sensor",
		logger.New("scanner"),
		func(entry discovery.Entry) {
			sub, err := remote.NewSuscriber(&entry)
			log.LogIfErr(err)
			if sub != nil {
				if _, present := sensors.Load(entry.Id); !present {
					log.Infof("new sensor found: %s", entry.Id)
					sensors.Store(entry.Id, sub)
					if len(sensorsIds()) == 1 {
						u.selectSensor(entry.Id)
					}
					wnd.Changed()
				}
			}
		},
		func(entry discovery.Entry) {
			log.Infof("removing sensor %s", entry.Id)
			sensors.Delete(entry.Id)
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

	//dev, _ = ak9753.NewMean(dev, 5)
	//defer dev.Close()
	//
	//detector = presence.New(dev, nil)
	//defer detector.Close()

	go u.updateSensorData()

	wnd.Main()
}

func sensorsIds() []string {
	ids := []string{}

	sensors.Range(func(key, value interface{}) bool {
		ids = append(ids, key.(string))
		return true
	})

	sort.Strings(ids)
	return ids
}
