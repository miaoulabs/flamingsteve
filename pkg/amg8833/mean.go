package amg8833

import (
	"flamingsteve/pkg/notification"
	"go.uber.org/atomic"
	"gonum.org/v1/gonum/stat"
	"sync"
)

type Mean struct {
	notification.NotifierImpl

	dev   Device
	state State

	sampleCount *atomic.Int32
	ir          [PIXEL_COUNT][]float64
	thermistor  []float64

	mutex sync.RWMutex
}

func NewMean(dev Device, samplesCount int) (*Mean, error) {

	if samplesCount < 2 {
		panic("sample count need 2+ for this to be useful")
	}

	if dev == nil {
		panic("cannnot use nil device")
	}

	m := &Mean{
		dev:         dev,
		sampleCount: atomic.NewInt32(int32(samplesCount)),
	}

	go m.run()

	return m, nil
}

func (m *Mean) SetSampleCount(count int) {
	if count < 2 {
		// todo: warn about invalid sample count
		count = 2
	}
	m.sampleCount.Store(int32(count))
}

func (m *Mean) run() {
	log.Infof("mean loop started")
	defer log.Infof("mean loop stopped")

	changed := make(chan bool)
	m.dev.Subscribe(changed)
	defer m.Unsubscribe(changed)

	for range changed {
		state := m.dev.State()

		m.mutex.Lock()

		pixels := state.Pixels

		max := int(m.sampleCount.Load())

		for i := range m.ir {
			m.ir[i] = append(m.ir[i], float64(pixels[i]))

			// drop the first item
			if len(m.ir[i]) > max {
				m.ir[i] = m.ir[i][len(m.ir[i])-max:]
			}
		}

		m.thermistor = append(m.thermistor, float64(state.Thermistor))
		if len(m.thermistor) > max {
			m.thermistor = m.thermistor[len(m.thermistor)-max:]
		}

		m.mutex.Unlock()
		m.Notify()
	}
}

func (m *Mean) Close() {
	m.UnsubscribeAll()
	m.dev.Close()
}

func (m *Mean) Thermistor() float32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return float32(stat.Mean(m.thermistor, nil))
}

func (m *Mean) Temperature(x, y int) float32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return float32(stat.Mean(m.ir[XYtoIndex(x, y)], nil))
}

func (m *Mean) Temperatures() [PIXEL_COUNT]float32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tmp := [PIXEL_COUNT]float32{}
	for i := range tmp {
		tmp[i] = float32(stat.Mean(m.ir[i], nil))
	}
	return tmp
}

func (m *Mean) State() State {
	return State{
		Pixels:     m.Temperatures(),
		Thermistor: m.Thermistor(),
	}
}
