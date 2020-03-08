package presence

import (
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/notification"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/atomic"
	"sync"
)

type Subscriber struct {
	notification.NotifierImpl
	sub     *nats.Subscription
	present *atomic.Bool
	sync.RWMutex

	entry discovery.Entry
}

func NewSubscriber(entry discovery.Entry) (*Subscriber, error) {
	s := &Subscriber{
		entry:   entry,
		present: atomic.NewBool(false),
	}

	if entry.Type != discovery.Detector {
		return nil, fmt.Errorf("wrong type of entry: expecting %s, has %d", discovery.Detector, entry.Type)
	}

	var err error
	s.sub, err = muthur.EncodedConnection().Subscribe(entry.DataTopic, s.update)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Subscriber) Close() {
	_ = s.sub.Unsubscribe()
}

func (s *Subscriber) IsPresent() bool {
	return s.present.Load()
}

func (s *Subscriber) Configs() []byte {
	data, err := s.entry.FetchConfig()
	if err != nil {
		log.Errorf("failed to fetch config from %s: %v", s.entry.Name, err)
	}
	return data
}

func (s *Subscriber) SetConfigs(data []byte) {
	err := s.entry.PushConfig(data)
	if err != nil {
		log.Errorf("failed to push config to %v: %v", s.entry.Name, err)
	}
}

func (s *Subscriber) update(state *DetectorState) {
	s.Lock()
	defer s.Unlock()
	if s.present.Load() != state.Present {
		log.Infof("remote detector state changed: %v", state.Present)
		s.present.Store(state.Present)
		s.Notify()
	}
}
