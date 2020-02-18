package display

import (
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"github.com/nats-io/nats.go"
)

type Listener struct {
	ident *discovery.Component
	dataConn *nats.Subscription
}

type DrawFunc func(msg *Message)

func NewListener(name string, model string, draw DrawFunc) *Listener {
	l := &Listener{}

	l.ident = discovery.NewComponent(discovery.IdentifierConfig{
		Name:  name,
		Model: model,
		Type:  discovery.Display,
	})
	l.ident.Connect()

	topic := l.ident.DataTopic()
	l.dataConn, _ = muthur.EncodedConnection().Subscribe(topic, draw)

	return l
}

func (l *Listener) Close() {
	_ = l.dataConn.Unsubscribe()
	l.ident.Close()
}
