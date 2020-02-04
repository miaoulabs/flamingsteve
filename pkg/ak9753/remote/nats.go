package remote

import (
	"github.com/nats-io/nats.go"
)

func natsErrorHandler(c *nats.Conn, subs *nats.Subscription, err error) {
	log.Errorf("error occurred on subs '%s': %w", subs.Subject, err)
}

func natsCloseHandler(conn *nats.Conn) {
	log.Infof("close connection to server %s", conn.ConnectedServerId())
}
