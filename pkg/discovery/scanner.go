package discovery

import (
	"flamingsteve/pkg/logger"
	"flamingsteve/pkg/muthur"
	"fmt"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

type Scanner struct {
	subOn  *nats.Subscription
	subOff *nats.Subscription
	log    logger.Logger
	tick   *time.Ticker

	mutex sync.Mutex
}

type OnHandler func(entry Entry)
type OffHandler func(entry Entry)

func NewScanner(etype EntryType, on OnHandler, off OffHandler) *Scanner {
	s := &Scanner{
		log: logFactory("scanner-" + string(etype)),
	}

	s.log.Infof("starting discovery for %s", etype)

	if on != nil {
		topic := fmt.Sprintf("%s.%s.>", TopicDeviceOn, etype)
		s.log.Infof("listening for topic '%s'", topic)
		s.subOn, _ = muthur.EncodedConnection().Subscribe(topic, on)
	}

	if off != nil {
		topic := fmt.Sprintf("%s.%s.>", TopicDeviceOff, etype)
		s.log.Infof("listening for topic '%s'", topic)
		s.subOff, _ = muthur.EncodedConnection().Subscribe(topic, off)
	}

	return s
}

func (s *Scanner) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.log.Infof("closing scanner")
	if s.subOff != nil {
		s.subOff.Unsubscribe()
	}

	if s.subOn != nil {
		s.subOn.Unsubscribe()
	}
}

func (s *Scanner) StartScan(interval time.Duration) {
	go func() {
		s.log.Infof("scan routine every %v ms", interval.Milliseconds())
		defer s.log.Infof("discovery scan stopped")

		s.mutex.Lock()
		s.tick = time.NewTicker(interval)
		s.mutex.Unlock()

		for range s.tick.C {
			s.Scan()
		}

		s.mutex.Lock()
		s.tick = nil
		s.mutex.Unlock()
	}()
}

func (s *Scanner) StopScan() {
	if s.tick != nil {
		s.tick.Stop()
	}
}

/*
	Send a broadcast message telling to tell all device to tell to
	broadcast their existence
*/
func (s *Scanner) Scan() {
	//s.log.Infof("broadcasting scan message on topic '%s'", TopicScan)
	err := muthur.Connection().Publish(TopicScan, []byte{})
	if err != nil {
		s.log.Errorf("error scanning: %v", err)
	}
}
