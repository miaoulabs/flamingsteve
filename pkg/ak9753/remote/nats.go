package remote

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

func natsErrorHandler(c *nats.Conn, subs *nats.Subscription, err error) {
	fmt.Fprintf(os.Stderr, "error occurred on subs '%s': %w\n", subs.Subject, err)
}

func natsCloseHandler(conn *nats.Conn) {
	fmt.Printf("close connection to server %s\n", conn.ConnectedServerId())
}

