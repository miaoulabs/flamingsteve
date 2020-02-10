package remote

import (
	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/notification"
	"github.com/nats-io/nats.go"
	"sync"
)

type Subscriber struct {
	sub   *nats.Subscription
	state ak9753.State
	mutex sync.RWMutex
	notification.NotifierImpl
}

func NewSuscriber(entry *discovery.Entry) (*Subscriber, error) {
	s := &Subscriber{}
	var err error

	if entry != nil {
		err = s.Change(*entry)
	}

	return s, err
}

func (s *Subscriber) Close() {
	log.Infof("closing subscriber")
	s.sub.Unsubscribe()
	s.UnsubscribeAll()
}

func (s *Subscriber) Change(entry discovery.Entry) error {
	s.sub.Unsubscribe()
	var err error
	s.sub, err = muthur.Connection().Subscribe(entry.DataTopic, s.update)
	return err
}

func (s *Subscriber) update(state *ak9753.State) {
	s.mutex.Lock()
	haschanged := !s.state.Equal(*state)
	s.state = *state
	s.mutex.Unlock()

	if haschanged {
		s.Notify()
	}
}

func (s *Subscriber) DeviceId() (uint8, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.DeviceId, nil
}

func (s *Subscriber) CompagnyCode() (uint8, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.CompagnyCode, nil
}

func (s *Subscriber) IR(idx int) (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Irs()[idx], nil
}

func (s *Subscriber) IR1() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Ir1, nil
}

func (s *Subscriber) IR2() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Ir2, nil
}

func (s *Subscriber) IR3() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Ir3, nil
}

func (s *Subscriber) IR4() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Ir4, nil
}

func (s *Subscriber) Temperature() (float32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Temperature, nil
}

func (s *Subscriber) All() ak9753.State {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state
}
