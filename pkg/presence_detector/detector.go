package pdetect

import (
	"flamingsteve/pkg/ak9753"
	"fmt"
	"os"
	"sync"
	"time"
)

type Options struct {
	Interval          time.Duration
	PresenceThreshold float32
	MovementThreshold float32
	Smoothing         float32
}

type Detector struct {
	dev  *ak9753.Device
	opts Options

	presence [FieldCount]bool
	movement uint8

	values [6]float32
	ders   [FieldCount]float32
	ders13 float32
	ders24 float32

	temp float32

	smoothers [smoothingCount]*smoother

	lastEval time.Time

	close chan bool

	mutex sync.RWMutex
}

func New(device *ak9753.Device, options *Options) *Detector {
	if options == nil {
		options = &Options{
			Interval:          time.Millisecond * 30,
			PresenceThreshold: 10,
			MovementThreshold: 10,
		}
	}

	d := &Detector{
		dev:   device,
		opts:  *options,
		close: make(chan bool),
	}

	for i := range d.smoothers {
		d.smoothers[i] = &smoother{
			avgWeigth: options.Smoothing, //0.3 very steep, 0.1 less steep, 0.05 less steep
		}
	}

	go d.run()

	return d
}

func (d *Detector) Close() {
	d.close <- true
	close(d.close)
}

func (d *Detector) PresentInField(idx int) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	r := d.presence[idx]
	//d.presence[idx] = false
	return r
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

func (d *Detector) PresenceAnyFields(clear bool) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	r := false
	for _, p := range d.presence {
		r = r || p
	}

	//fmt.Printf("p: %v\n", d.presence)

	if clear { // reset presence
		d.presence = [FieldCount]bool{}
	}

	return r
}

func (d *Detector) Temperature() float32 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.temp
}

func (d *Detector) IR1() float32 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.smoothers[0].last
}

func (d *Detector) IR2() float32 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.smoothers[1].last
}

func (d *Detector) IR3() float32 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.smoothers[2].last
}

func (d *Detector) IR4() float32 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.smoothers[3].last
}

func (d *Detector) DerivativeOfIR(idx int) float32 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.ders[idx]
}

func (d *Detector) DerivativeOfIR1() float32 {
	return d.DerivativeOfIR(0)
}

func (d *Detector) DerivativeOfIR2() float32 {
	return d.DerivativeOfIR(1)
}

func (d *Detector) DerivativeOfIR3() float32 {
	return d.DerivativeOfIR(2)
}

func (d *Detector) DerivativeOfIR4() float32 {
	return d.DerivativeOfIR(3)
}

func (d *Detector) DerivativeOfDiff13() float32 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.ders13
}

func (d *Detector) DerivativeOfDiff24() float32 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.ders24
}

func (d *Detector) Movement() uint8 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	r := d.movement
	d.movement = MovementNone
	return r
}

func (d *Detector) run() {

	println("starting detection loop")
	defer println("detection loop stopped")

	tick := time.NewTicker(time.Millisecond * 5)
	defer tick.Stop()

	for range tick.C {
		select {
		case <-d.close:
			return
		default:
		}

		if !d.dev.DataReady() {
			continue
		}

		ir1, err := d.dev.IR1()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample: %w", err)
			continue
		}

		ir2, err := d.dev.IR2()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample: %w", err)
			continue
		}

		ir3, err := d.dev.IR3()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample: %w", err)
			continue
		}

		ir4, err := d.dev.IR4()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample: %w", err)
			continue
		}

		temp, err := d.dev.Temperature()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading temperature: %w", err)
		}

		err = d.dev.StartNextSample()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample: %w", err)
			continue
		}

		diff13 := ir1 - ir3
		diff24 := ir2 - ir4

		//fmt.Printf("Reading: IR1: %f, IR2: %f, IR3: %f, IR4: %f, Diff13: %f, Diff24: %f\n",
		//	ir1, ir2, ir3, ir4, diff13, diff24,
		//)

		d.mutex.Lock()
		d.temp = temp
		d.smoothers[0].add(ir1)
		d.smoothers[1].add(ir2)
		d.smoothers[2].add(ir3)
		d.smoothers[3].add(ir4)
		d.smoothers[4].add(diff13)
		d.smoothers[5].add(diff24)

		if time.Now().Sub(d.lastEval) > d.opts.Interval {
			for i := 0; i < FieldCount; i++ {
				der := d.smoothers[i].derivative()
				d.ders[i] = der

				//fmt.Printf("d#%d: %f ", i, der)

				if der > d.opts.PresenceThreshold {
					d.presence[i] = true
				} else if der < -d.opts.PresenceThreshold {
					d.presence[i] = false
				}
			}

			//println()

			d.ders13 = d.smoothers[4].derivative()
			if d.ders13 > d.opts.PresenceThreshold {
				d.movement &= 0b11111100
				d.movement |= MovementFrom3to1
			} else if d.ders13 < -d.opts.PresenceThreshold {
				d.movement &= 0b11111100
				d.movement |= MovementFrom1to3
			}

			d.ders24 = d.smoothers[5].derivative()
			if d.ders24 > d.opts.PresenceThreshold {
				d.movement &= 0b11110011
				d.movement |= MovementFrom4to2
			} else if d.ders24 < -d.opts.PresenceThreshold {
				d.movement &= 0b11110011
				d.movement |= MovementFrom2to4
			}

			d.lastEval = time.Now()
		}

		d.mutex.Unlock()
	}
}
