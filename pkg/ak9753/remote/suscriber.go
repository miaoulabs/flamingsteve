package remote

import (
	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/notify"
	"github.com/nats-io/nats.go"
	"sync"
)

type Suscriber struct {
	conn    *nats.Conn
	encoded *nats.EncodedConn
	subs    *nats.Subscription

	state ak9753.State
	mutex sync.RWMutex

	notify.Notifier
}

func NewSuscriber(url string, name string) (*Suscriber, error) {
	s := &Suscriber{}

	if url == "" {
		url = "nats//localhost:4222"
	}

	if name == "" {
		name = "ak9753-suscriber"
	}

	var err error
	s.conn, err = nats.Connect(url,
		nats.Name(name),
		nats.ErrorHandler(natsErrorHandler),
		nats.ClosedHandler(natsCloseHandler),
	)
	if err != nil {
		return nil, err
	}

	s.encoded, err = nats.NewEncodedConn(s.conn, nats.JSON_ENCODER)
	if err != nil {
		s.conn.Close()
		return nil, err
	}

	s.subs, err = s.encoded.Subscribe(Topic, func(state *ak9753.State) {
		s.mutex.Lock()
		haschanged := !s.state.Equal(*state)
		s.state = *state
		s.mutex.Unlock()

		if haschanged {
			s.Notify()
		}
	})

	if err != nil {
		s.encoded.Close()
	}

	return s, nil
}

func (s *Suscriber) Close() {
	log.Infof("closing suscriber")
	s.encoded.Close()
	s.UnsubscribeAll()
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
