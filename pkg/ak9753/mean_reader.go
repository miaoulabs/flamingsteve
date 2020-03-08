package ak9753

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
	ir          [FieldCount][]float64
	mutex       sync.RWMutex
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
		state := m.dev.All()

		irs := state.Irs()

		for i := range m.ir {
			m.mutex.Lock()
			m.ir[i] = append(m.ir[i], float64(irs[i]))

			// drop the first item
			max := int(m.sampleCount.Load())
			if len(m.ir[i]) > max {
				m.ir[i] = m.ir[i][len(m.ir[i])-max:]
			}
			m.mutex.Unlock()
		}

		m.Notify()
	}
}

func (m *Mean) Close() {
	m.UnsubscribeAll()
	m.dev.Close()
}

func (m *Mean) DeviceId() (uint8, error) {
	return m.dev.DeviceId()
}

func (m *Mean) CompagnyCode() (uint8, error) {
	return m.dev.CompagnyCode()
}

func (m *Mean) IR(idx int) (float32, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return float32(stat.Mean(m.ir[idx], nil)), nil
}

func (m *Mean) IR1() (float32, error) {
	return m.IR(0)
}

func (m *Mean) IR2() (float32, error) {
	return m.IR(1)
}

func (m *Mean) IR3() (float32, error) {
	return m.IR(2)
}

func (m *Mean) IR4() (float32, error) {
	return m.IR(3)
}

func (m *Mean) Temperature() (float32, error) {
	return m.dev.Temperature()
}

func (m *Mean) All() State {
	st := m.dev.All()
	st.Ir1, _ = m.IR(0)
	st.Ir2, _ = m.IR(1)
	st.Ir3, _ = m.IR(2)
	st.Ir4, _ = m.IR(3)
	return st
}
