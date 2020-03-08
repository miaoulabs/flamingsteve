package display

import (
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
)

type Listener struct {
	ident    *discovery.Component
	dataConn *nats.Subscription
	info     Info
}

type DrawFunc func(msg *Message)

func NewListener(name string, model string, width, height uint, draw DrawFunc) *Listener {
	l := &Listener{
		info: Info{
			Width:  width,
			Height: height,
		},
	}

	l.ident = discovery.NewComponent(discovery.IdentifierConfig{
		Name:  name,
		Model: model,
		Type:  discovery.Display,
	})
	l.ident.OnConfigRequest(l.sendInfo)
	l.ident.Connect()

	topic := l.ident.DataTopic()
	l.dataConn, _ = muthur.EncodedConnection().Subscribe(topic, draw)

	return l
}

func (l *Listener) Close() {
	_ = l.dataConn.Unsubscribe()
	l.ident.Close()
}

func (l *Listener) sendInfo() []byte {
	data, _ := jsoniter.Marshal(l.info)
	// todo: log error ?
	return data
}
