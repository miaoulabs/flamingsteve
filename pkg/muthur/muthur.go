package muthur

import (
	"errors"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

var (
	once     sync.Once
	natsConn *nats.Conn
	natsEnc  *nats.EncodedConn
)

func Connect(clientName string) {
	once.Do(func() {
		dns := MustResolveServer()

		err := errors.New("")
		for err != nil {
			natsConn, err = nats.Connect(
				dns.HostName,
				nats.Name(clientName),
				nats.ErrorHandler(natsErrorHandler),
				nats.ClosedHandler(natsCloseHandler),
			)
			if err != nil {
				time.Sleep(time.Millisecond * 200)
			}
		}

		natsEnc, _ = nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	})
}

func Connection() *nats.EncodedConn {
	if natsConn == nil {
		panic("muthur.Connect() has not been called")
	}
	return natsEnc
}

func Close() {
	natsEnc.Close()
}

func natsErrorHandler(c *nats.Conn, subs *nats.Subscription, err error) {
	log.Errorf("error occurred on subs '%s': %v", subs.Subject, err)
}

func natsCloseHandler(conn *nats.Conn) {
	log.Infof("close connection to server %s", conn.ConnectedServerId())
}
