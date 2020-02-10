package presence

import (
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/notification"
)

type Publisher struct {
	notification.NotifierImpl
	detector Detector
	ident *discovery.Component
	changed chan bool
}

func (p *Publisher) Configs() []byte {
	panic("implement me")
}

func (p *Publisher) SetConfigs(data []byte) {
	panic("implement me")
}

func NewPublisher(detector Detector, ident *discovery.Component) *Publisher {
	p := &Publisher{
		detector: detector,
		ident: ident,
		changed: make(chan bool),
	}

	ident.Connect()

	go p.run()
	return p
}

func (p *Publisher) Close() {
	p.detector.Unsubscribe(p.changed)
	close(p.changed)
}

func (p *Publisher) IsPresent() bool {
	return p.detector.IsPresent()
}

func (p *Publisher) run() {
	log.Infof("publisher loop started")
	defer log.Infof("publisher loop stopped")

	p.detector.Subscribe(p.changed)

	for range p.changed {

		// todo: add log
		err := p.ident.PushData(&DetectorState{
			Present: p.IsPresent(),
		})
		if err != nil {
			log.Errorf("error sending data: %v", err)
		}

		p.Notify()
	}
}
