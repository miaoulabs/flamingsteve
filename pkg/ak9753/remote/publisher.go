package remote

import (
	"flamingsteve/pkg/ak9753"
	"github.com/nats-io/nats.go"
	"time"
)

/*
	This will connect to a ak9753 sensor and output all it's data
	into a nats topics.
*/
type Publisher struct {
	dev     ak9753.Device
	conn    *nats.Conn
	encoded *nats.EncodedConn
	close   chan bool
}

func NewPublisher(dev ak9753.Device, url string) (*Publisher, error) {
	if dev == nil {
		panic("ak9753 cannot be nil")
	}

	p := &Publisher{
		dev:   dev,
		close: make(chan bool),
	}

	if url == "" {
		url = "nats://localhost:4222"
	}

	var err error
	p.conn, err = nats.Connect(url,
		nats.Name("ak9753"),
		nats.ErrorHandler(natsErrorHandler),
		nats.ClosedHandler(natsCloseHandler),
	)
	if err != nil {
		return nil, err
	}

	p.encoded, err = nats.NewEncodedConn(p.conn, nats.JSON_ENCODER)
	if err != nil {
		p.conn.Close()
		return nil, err
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
	defer p.encoded.Close()
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
			err := p.encoded.Publish(Topic, last)
			if err != nil {
				log.Errorf("error sending message: %w", err)
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

func (p *Publisher) Subscribe(channel chan<- bool) {
	p.dev.Subscribe(channel)
}
