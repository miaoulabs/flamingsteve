package remote

import (
	"flamingsteve/pkg/amg8833"
	"flamingsteve/pkg/discovery"
	"github.com/nats-io/nats.go"
)

/*
	This will connect to a ak9753 sensor and output all it's data
	into a nats topics.
*/
type Publisher struct {
	dev   amg8833.Device
	conn  *nats.Conn
	ident *discovery.Component
}

func NewPublisher(dev amg8833.Device, ident *discovery.Component) (*Publisher, error) {
	if dev == nil {
		panic("amg8833 cannot be nil")
	}

	p := &Publisher{
		dev:   dev,
		ident: ident,
	}

	go p.run()

	return p, nil
}

func (p *Publisher) Close() {
	p.dev.Close()
	log.Infof("publisher closed")
}

func (p *Publisher) run() {
	log.Infof("publisher loop started")
	defer log.Infof("publisher loop stopped")

	changed := make(chan bool)
	p.dev.Subscribe(changed)
	defer p.dev.Unsubscribe(changed)

	last := amg8833.State{}

	for range changed {
		newV := p.dev.State()

		if !last.Equal(newV) {
			err := p.ident.PushData(newV)
			if err != nil {
				log.Errorf("error sending message: %v", err)
			}
		}
	}
}

func (p *Publisher) Thermistor() float32 {
	return p.dev.Thermistor()
}

func (p *Publisher) Temperature(x, y int) float32 {
	return p.dev.Temperature(x, y)
}

func (p *Publisher) Temperatures() [amg8833.PIXEL_COUNT]float32 {
	return p.dev.Temperatures()
}

func (p *Publisher) State() amg8833.State {
	return p.dev.State()
}

func (p *Publisher) Subscribe(channel chan<- bool) {
	p.dev.Subscribe(channel)
}

func (p *Publisher) Unsubscribe(channel chan<- bool) {
	p.dev.Unsubscribe(channel)
}
