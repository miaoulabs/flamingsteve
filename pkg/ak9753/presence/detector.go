package presence

import (
	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/notification"
	"sync"
	"time"
)

type Detector struct {
	notification.NotifierImpl
	dev     ak9753.Device
	mutex   sync.RWMutex
	opts    Options
	changed chan bool
	mean    *ak9753.Mean

	since [ak9753.FieldCount]*time.Time
}

func New(device ak9753.Device, options *Options) (*Detector, error) {
	if options == nil {
		options = &Options{
			Delay:          DefaultDelay,
			Smoothing:      DefaultSmoothing,
			MinimumSensors: DefaultMinSensor,
			Threshold:      DefaultThreshold,
		}
	}

	mean, err := ak9753.NewMean(device, options.Smoothing)
	if err != nil {
		return nil, err
	}

	d := &Detector{
		dev:     device,
		changed: make(chan bool),
		mean: mean,
		since: [ak9753.FieldCount]*time.Time{},
	}

	d.ApplyOptions(*options)

	go d.run()

	return d, nil
}

func (d *Detector) Close() {
	log.Infof("closing detector")
	d.dev.Unsubscribe(d.changed)
}

func (d *Detector) Options() Options {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.opts
}


func (d *Detector) Configs() []byte {
	return d.opts.Marshal()
}

func (d *Detector) SetConfigs(data []byte) {
	opts := UnmarshalOptions(data)
	d.ApplyOptions(opts)
}

func (d *Detector) ApplyOptions(opts Options) {

	if opts.MinimumSensors <= 0 {
		// todo: log warning about invalid sensor count
		opts.MinimumSensors = 1
	}

	d.mean.SetSampleCount(opts.MinimumSensors)

	d.mutex.Lock()
	d.opts = opts
	d.mutex.Unlock()
}

func (d *Detector) IsPresent() bool {
	count := 0

	opts := d.Options() // thread safe local copy

	for idx := 0; idx < ak9753.FieldCount; idx++ {
		if d.PresentInField(idx) {
			count++
		}
	}

	return count >= opts.MinimumSensors
}

func (d *Detector) PresentInField(idx int) bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	delay := d.since[idx]

	if delay == nil {
		return false
	} else {
		return time.Now().Sub(*delay) > d.opts.Delay
	}
}

func (d *Detector) PresentInField1() bool {
	return d.PresentInField(0)
}

func (d *Detector) PresentInField2() bool {
	return d.PresentInField(1)
}

func (d *Detector) PresentInField3() bool {
	return d.PresentInField(2)
}

func (d *Detector) PresentInField4() bool {
	return d.PresentInField(3)
}

func (d *Detector) run() {
	d.dev.Subscribe(d.changed)

	for range d.changed {
		d.mutex.Lock()

		for idx := range d.since {
			val, _ := d.mean.IR(idx)

			if val > d.opts.Threshold {
				if d.since[idx] == nil {
					now := time.Now()
					d.since[idx] = &now
				}
			} else {
				d.since[idx] = nil
			}
		}

		d.mutex.Unlock()
	}
}
