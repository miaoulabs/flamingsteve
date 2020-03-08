package remote

import (
	"github.com/nats-io/nats.go"
	"sync"

	"flamingsteve/pkg/amg8833"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/notification"
)

type Subscriber struct {
	sub   *nats.Subscription
	state amg8833.State
	mutex sync.RWMutex
	notification.NotifierImpl
}

func NewSuscriber(entry discovery.Entry) (*Subscriber, error) {
	s := &Subscriber{}
	var err error
	err = s.Change(entry)
	return s, err
}

func (s *Subscriber) Close() {
	log.Infof("closing subscriber")
	s.sub.Unsubscribe()
	s.UnsubscribeAll()
}

func (s *Subscriber) Thermistor() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Thermistor
}

func (s *Subscriber) Temperature(x, y int) float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Pixels[x+y*amg8833.ROW_COUNT]
}

func (s *Subscriber) Temperatures() [amg8833.PIXEL_COUNT]float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state.Pixels
}

func (s *Subscriber) State() amg8833.State {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state
}

func (s *Subscriber) Change(entry discovery.Entry) error {
	_ = s.sub.Unsubscribe()
	var err error
	s.sub, err = muthur.EncodedConnection().Subscribe(entry.DataTopic, s.update)
	return err
}

func (s *Subscriber) update(state *amg8833.State) {
	haschanged := false

	s.mutex.Lock()
	if state != nil {
		haschanged = !s.state.Equal(*state)
		s.state = *state
	}
	s.mutex.Unlock()

	if haschanged {
		s.Notify()
	}
}
