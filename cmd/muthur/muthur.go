package main

import (
	"flamingsteve/pkg/muthur"
	"github.com/draeron/gopkgs/logger"
	"github.com/grandcat/zeroconf"
	natsd "github.com/nats-io/nats-server/v2/server"
	"os"
	"os/signal"
	"time"
)

func main() {
	log := logger.New("main")

	log.Info("started")
	defer log.Info("stopped")

	if alreadyExists() {
		log.Errorf("another muthur exists, THERE CAN BE ONLY ONE!")
		return
	}

	log.Info("registering zeroconf dns")
	mdns, err := zeroconf.Register("muthur",
		muthur.ZeroConfServiceName,
		muthur.ZeroConfDomain,
		4222,
		[]string{},
		muthur.ListMulticastInterfaces(),
	)
	log.StopIfErr(err)
	defer mdns.Shutdown()

	opts := &natsd.Options{
		ServerName: "muthur",
		HTTPPort:   8222,
		Debug:      true,
	}
	svr, err := natsd.NewServer(opts)
	log.StopIfErr(err)

	svr.SetLogger(&natsLogger{log}, true, false)

	go func() {
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs)

		go func() {
			sig := <-sigs
			log.Infof("received signal %d, stopping application", sig)
			done <- true
		}()

		<-done

		svr.Shutdown()
	}()

	log.Info("starting nats server")
	svr.Start()
}

func alreadyExists() bool {
	others := muthur.ResolveServers(time.Second * 3)
	return len(others) > 0
}
