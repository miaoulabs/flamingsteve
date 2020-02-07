package remote

import (
	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/discovery"
	"github.com/nats-io/nats.go"
	"time"
)

/*
	This will connect to a ak9753 sensor and output all it's data
	into a nats topics.
*/
type Publisher struct {
	dev   ak9753.Device
	conn  *nats.Conn
	close chan bool
	ident *discovery.Identifier
}

func NewPublisher(dev ak9753.Device, ident *discovery.Identifier) (*Publisher, error) {
	if dev == nil {
		panic("ak9753 cannot be nil")
	}

	p := &Publisher{
		dev:   dev,
		close: make(chan bool),
		ident: ident,
	}

	go p.run()

	return p, nil
}

func (p *Publisher) Close() {
	p.dev.Close()
	log.Infof("closing publisher")
	p.close <- true
	close(p.close)
}

func (p *Publisher) run() {
	log.Infof("publisher loop started")
	defer log.Infof("publisher loop stopped")

	last := ak9753.State{}

	tick := time.NewTicker(time.Millisecond * 5)
	defer tick.Stop()

	for range tick.C {
		select {
		case <-p.close:
			return // exit loop
		default:
		}

		newV := p.dev.All()

		if !newV.Equal(last) {
			last = newV
			err := p.ident.PushRaw(last)
			if err != nil {
				log.Errorf("error sending message: %v", err)
			}
		}
	}
}

func (p *Publisher) DeviceId() (uint8, error) {
	return p.dev.DeviceId()
}

func (p *Publisher) CompagnyCode() (uint8, error) {
	return p.dev.CompagnyCode()
}

func (p *Publisher) IR(idx int) (float32, error) {
	return p.dev.IR(idx)
}

func (p *Publisher) IR1() (float32, error) {
	return p.dev.IR1()
}

func (p *Publisher) IR2() (float32, error) {
	return p.dev.IR2()
}

func (p *Publisher) IR3() (float32, error) {
	return p.dev.IR3()
}

func (p *Publisher) IR4() (float32, error) {
	return p.dev.IR4()
}

func (p *Publisher) Temperature() (float32, error) {
	return p.dev.Temperature()
}

func (p *Publisher) All() ak9753.State {
	return p.dev.All()
}

func (p *Publisher) UnsubscribeAll() {
	p.dev.UnsubscribeAll()
}

func (p *Publisher) Subscribe(channel chan<- bool) {
	p.dev.Subscribe(channel)
}
