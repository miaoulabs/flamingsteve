package discovery

import (
	"flamingsteve/pkg/logger"
	"flamingsteve/pkg/muthur"
	"fmt"
	"github.com/nats-io/nats.go"
)

type Scanner struct {
	subOn  *nats.Subscription
	subOff *nats.Subscription
	log    logger.Logger
}

type OnHandler func(entry Entry)
type OffHandler func(entry Entry)

func NewScanner(etype EntryType, log logger.Logger, on OnHandler, off OffHandler) *Scanner {
	s := &Scanner{
		log: logger.Dummy(),
	}

	if log != nil {
		s.log = log
	}

	s.log.Infof("starting discovery for %s", etype)

	if on != nil {
		topic := fmt.Sprintf("%s.%s.>", TopicDeviceOn, etype)
		s.log.Infof("listening for topic '%s'", topic)
		s.subOn, _ = muthur.Connection().Subscribe(topic, on)
	}

	if off != nil {
		topic := fmt.Sprintf("%s.%s.>", TopicDeviceOff, etype)
		s.log.Infof("listening for topic '%s'", topic)
		s.subOff, _ = muthur.Connection().Subscribe(topic, off)
	}

	return s
}

func (s *Scanner) Close() {
	s.log.Infof("closing scanner")
	if s.subOff != nil {
		s.subOff.Unsubscribe()
	}

	if s.subOn != nil {
		s.subOn.Unsubscribe()
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
