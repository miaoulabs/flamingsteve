package ak9753

import (
	"flamingsteve/pkg/notify"
	"fmt"
	"os"
	"sync"
	"time"
)

/*
	Thread safe implementation which read from a physical ak9753 device
	and store it's data
*/
type Reader struct {
	dev     *Physical
	close   chan bool
	mutex   sync.RWMutex
	state   State
	lastErr error

	notify.Notifier
}

func NewReader(dev *Physical) (*Reader, error) {
	r := &Reader{
		dev: dev,
		close: make(chan bool),
	}

	var err error

	err = r.initDevice()
	if err != nil {
		return nil, err
	}

	go r.run()

	return r, nil
}

func (r *Reader) Close() {
	println("closing ak9753 reader")
	r.close <- true
	close(r.close)
	r.UnsubscribeAll()
}

func (r *Reader) initDevice() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var err error

	r.state.CompagnyCode, err = r.dev.CompagnyCode()
	if err != nil {
		return err
	}

	r.state.DeviceId, err = r.dev.DeviceId()
	if err != nil {
		return err
	}

	model, err := r.dev.Model()
	if err != nil {
		return err
	}
	fmt.Printf("sensor model: %s\n", model)

	return nil
}

func (r *Reader) run() {
	println("ak9753 reader loop started")
	defer println("ak9753 reader loop stopped")

	var err error

	err = r.dev.StartNextSample()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting next sample: %w", err)
	}

	tick := time.NewTicker(time.Millisecond * 5)
	defer tick.Stop()

	for range tick.C {
		select {
		case <-r.close:
			return // exit loop
		default:
		}

		if !r.dev.DataReady() {
			continue
		}

		var err error

		state := State{}

		state.Temperature, err = r.dev.Temperature()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading temperature: %w", err)
		}

		state.Ir1, err = r.dev.IR1()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample for ir1: %w", err)
		}

		state.Ir2, err = r.dev.IR2()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample: %w", err)
		}

		state.Ir3, err = r.dev.IR3()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample: %w", err)
		}

		state.Ir4, err = r.dev.IR4()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading sample: %w", err)
		}

		err = r.dev.StartNextSample()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error starting next sample: %w", err)
		}

		// update state
		r.mutex.Lock()
		haschanged := !r.state.Equal(state)
		r.state = state
		r.mutex.Unlock()

		if haschanged {
			r.Notify()
		}
	}
}

func (r *Reader) DeviceId() (uint8, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.DeviceId, nil
}

func (r *Reader) CompagnyCode() (uint8, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.CompagnyCode, nil
}

func (r *Reader) IR1() (float32, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.Ir1, nil
}

func (r *Reader) IR2() (float32, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.Ir2, nil
}

func (r *Reader) IR3() (float32, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.Ir3, nil
}

func (r *Reader) IR4() (float32, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.Ir4, nil
}

func (r *Reader) Temperature() (float32, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state.Temperature, nil
}

func (r *Reader) All() State {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state
}
