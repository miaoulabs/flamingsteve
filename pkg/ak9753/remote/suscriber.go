package remote

import (
	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/notify"
	"github.com/nats-io/nats.go"
	"sync"
)

type Suscriber struct {
	sub   *nats.Subscription
	state ak9753.State
	mutex sync.RWMutex
	notify.Notifier
}

func NewSuscriber(entry *discovery.Entry) (*Suscriber, error) {
	s := &Suscriber{}
	var err error

	if entry != nil {
		err = s.Change(*entry)
	}

	return s, err
}

func (s *Suscriber) Close() {
	log.Infof("closing suscriber")
	s.sub.Unsubscribe()
	s.UnsubscribeAll()
}

func (s *Suscriber) Change(entry discovery.Entry) error {
	s.sub.Unsubscribe()
	var err error
	s.sub, err = muthur.Connection().Subscribe(entry.DataTopic, s.update)
	return err
}

func (s *Suscriber) update(state *ak9753.State) {
	s.mutex.Lock()
	haschanged := !s.state.Equal(*state)
	s.state = *state
	s.mutex.Unlock()

	if haschanged {
		s.Notify()
	}
}

func (s *Suscriber) DeviceId() (uint8, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.DeviceId, nil
}

func (s *Suscriber) CompagnyCode() (uint8, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.CompagnyCode, nil
}

func (s *Suscriber) IR(idx int) (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Irs()[idx], nil
}

func (s *Suscriber) IR1() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Ir1, nil
}

func (s *Suscriber) IR2() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Ir2, nil
}

func (s *Suscriber) IR3() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Ir3, nil
}

func (s *Suscriber) IR4() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Ir4, nil
}

func (s *Suscriber) Temperature() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Temperature, nil
}

func (s *Suscriber) All() ak9753.State {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state
}
