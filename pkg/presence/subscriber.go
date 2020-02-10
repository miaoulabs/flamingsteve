package presence

import (
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/notification"
	"fmt"
	"github.com/nats-io/nats.go"
	"sync"
)

type Subscriber struct {
	notification.NotifierImpl
	sub   *nats.Subscription
	present bool
	sync.RWMutex
}

func NewSubscriber(entry *discovery.Entry) (*Subscriber, error) {
	s := &Subscriber{}

	if entry.Type != discovery.Detector {
		return nil, fmt.Errorf("wrong type of entry: expecting %s, has %d", discovery.Detector, entry.Type)
	}

	var err error
	s.sub, err = muthur.Connection().Subscribe(entry.DataTopic, s.update)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Subscriber) Close() {
	s.sub.Unsubscribe()
}

func (s *Subscriber) IsPresent() bool {
	s.RLock()
	defer s.RUnlock()
	return s.present
}

func (s *Subscriber) Configs() []byte {
	panic("implement me")
}

func (s *Subscriber) SetConfigs(data []byte) {
	panic("implement me")
}

func (s *Subscriber) update(state *DetectorState) {
	s.Lock()
	defer s.Unlock()
	if s.present != state.Present {
		log.Infof("remote detector state changed: %v", state.Present)
		s.present = state.Present
		s.Notify()
	}
}
