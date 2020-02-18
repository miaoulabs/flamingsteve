package presence

import (
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/notification"
)

type Publisher struct {
	notification.NotifierImpl
	detector  Detector
	component *discovery.Component
	changed   chan bool
}

func (p *Publisher) Configs() []byte {
	return p.detector.Configs()
}

func (p *Publisher) SetConfigs(data []byte) {
	p.detector.SetConfigs(data)
}

func NewPublisher(detector Detector, component *discovery.Component) *Publisher {
	p := &Publisher{
		detector:  detector,
		component: component,
		changed:   make(chan bool),
	}

	component.Connect()
	component.OnConfigWrite(p.SetConfigs)
	component.OnConfigRequest(p.Configs)

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
		err := p.component.PushData(&DetectorState{
			Present: p.IsPresent(),
		})
		if err != nil {
			log.Errorf("error sending data: %v", err)
		}

		p.Notify()
	}
}
