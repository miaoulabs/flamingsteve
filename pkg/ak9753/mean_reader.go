package ak9753

import (
	"flamingsteve/pkg/notify"
	"gonum.org/v1/gonum/stat"
)

type Mean struct {
	dev Device
	state State
	notify.Notifier

	sampleCount int
	ir          [FieldCount][]float64
}

func NewMean(dev Device, samplesCount int) (*Mean, error) {

	if samplesCount < 2 {
		panic("sample count need 2+ for this to be useful")
	}

	if dev == nil {
		panic("cannnot use nil device")
	}

	m := &Mean{
		dev: dev,
		sampleCount: samplesCount,
	}

	go m.run()

	return m, nil
}

func (m *Mean) run () {
	println("mean loop started")
	defer println("mean loop stopped")

	changed := make(chan bool)
	m.dev.Subscribe(changed)

	for range changed {
		state := m.dev.All()

		irs := state.Irs()

		for i := range m.ir {
			m.ir[i] = append(m.ir[i], float64(irs[i]))

			// drop the first item
			if len(m.ir[i]) > m.sampleCount {
				m.ir[i] = m.ir[i][1:]
			}
		}

		m.Notify()
	}
}

func (m *Mean) Close() {
	m.dev.Close()
	m.UnsubscribeAll()
}

func (m *Mean) DeviceId() (uint8, error) {
	return m.dev.DeviceId()
}

func (m *Mean) CompagnyCode() (uint8, error) {
	return  m.dev.CompagnyCode()
}

func (m *Mean) IR(idx int) float32 {
	return float32(stat.Mean(m.ir[0], nil))
}

func (m *Mean) IR1() (float32, error) {
	return m.IR(0), nil
}

func (m *Mean) IR2() (float32, error) {
	return m.IR(1), nil
}

func (m *Mean) IR3() (float32, error) {
	return m.IR(2), nil
}

func (m *Mean) IR4() (float32, error) {
	return m.IR(3), nil
}

func (m *Mean) Temperature() (float32, error) {
	return m.dev.Temperature()
}

func (m *Mean) All() State {
	st := m.dev.All()
	st.Ir1 = m.IR(0)
	st.Ir2 = m.IR(1)
	st.Ir3 = m.IR(2)
	st.Ir4 = m.IR(3)
	return st
}
