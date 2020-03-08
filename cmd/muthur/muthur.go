package main

import (
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"time"

	"flamingsteve/cmd"
	"flamingsteve/pkg/muthur"
	"github.com/draeron/gopkgs/logger"
	"github.com/grandcat/zeroconf"
	natsd "github.com/nats-io/nats-server/v2/server"
)

var args struct {
	natsPort int
	mgtPort  int
}

func init() {
	pflag.IntVarP(&args.natsPort, "nats-port", "p", natsd.DEFAULT_PORT, "port for nats server")
	pflag.IntVarP(&args.mgtPort, "mgt-port", "m", natsd.DEFAULT_HTTP_PORT, "port for nats server")
}

func main() {
	pflag.Parse()
	cmd.SetupLoggers()

	log := logger.New("main")

	log.Info("muthur started")
	defer log.Info("muthur  stopped")

	if alreadyExists() {
		log.Errorf("another muthur exists, THERE CAN BE ONLY ONE!")
		return
	}

	log.Info("registering zeroconf dns")
	mdns, err := zeroconf.Register("muthur",
		muthur.ZeroConfServiceName,
		muthur.ZeroConfDomain,
		args.natsPort,
		[]string{},
		muthur.ListMulticastInterfaces(),
	)
	for _, inet := range muthur.ListMulticastInterfaces() {
		log.Infof("listening for broadcast on interface %s", inet.Name)
	}
	log.StopIfErr(err)
	defer mdns.Shutdown()

	opts := &natsd.Options{
		ServerName: "muthur",
		HTTPPort:   args.mgtPort,
		Port:       args.natsPort,
		Debug:      true,
	}
	svr, err := natsd.NewServer(opts)
	log.StopIfErr(err)

	svr.SetLogger(&natsLogger{log}, true, false)

	go func() {
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, cmd.TerminationSignals...)

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
