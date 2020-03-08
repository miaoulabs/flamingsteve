package amg8833

import (
	"flamingsteve/pkg/notification"
	"sync"
	"time"
)

type Reader struct {
	notification.NotifierImpl

	dev   *Physical
	close chan bool
	mutex sync.RWMutex
	state State
}

func NewReader(device *Physical) (*Reader, error) {
	r := &Reader{
		dev:   device,
		close: make(chan bool),
	}

	go r.run()

	return r, nil
}

func (r *Reader) Thermistor() float32 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.Thermistor
}

func (r *Reader) Temperature(x, y int) float32 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.Pixel(x, y)
}

func (r *Reader) Temperatures() [PIXEL_COUNT]float32 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.Pixels
}

func (r *Reader) State() State {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state
}

func (r *Reader) Close() {
	r.close <- true
	close(r.close)
}

func (r *Reader) run() {
	log.Infof("amg8833 reader loop started")
	defer log.Infof("amg8833 reader loop stopped")

	tick := time.NewTicker(time.Second / 10)
	defer tick.Stop()

	for {
		select {
		case <-r.close:
			return // exit loop
		default:
		}

		var err error
		readState := State{}
		readState.Pixels, err = r.dev.PixelTemperature()
		if err != nil {
			log.Errorf("failed to read temperature from sensor: %s", err)
			continue
		}

		readState.Thermistor, err = r.dev.Thermistor()
		if err != nil {
			log.Errorf("failed to read thermistor from sensor: %s", err)
			continue
		}

		r.mutex.Lock()
		changed := !r.state.Equal(readState)
		r.state = readState
		r.mutex.Unlock()

		if changed {
			r.Notify()
		}
	}
}
