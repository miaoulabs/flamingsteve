package main

import (
	"flamingsteve/pkg/discovery"
	"github.com/draeron/gopkgs/logger"
	"github.com/grandcat/zeroconf"
	natsd "github.com/nats-io/nats-server/v2/server"
	"net"
	"os"
	"os/signal"
	"strings"
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
		discovery.ZeroConfServiceName,
		discovery.ZeroConfDomain,
		4222,
		[]string{},
		listMulticastInterfaces(),
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

/*
	This a copy/paste from zeroconf package, except we filter
	vEthernet (docker) interface
 */
func listMulticastInterfaces() []net.Interface {
	var interfaces []net.Interface
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, ifi := range ifaces {
		if (ifi.Flags & net.FlagUp) == 0 {
			continue
		}
		if strings.Contains(ifi.Name, "vEthernet") {
			continue
		}

		if (ifi.Flags & net.FlagMulticast) > 0 {
			interfaces = append(interfaces, ifi)
		}
	}

	return interfaces
}

func alreadyExists() bool {
	others := discovery.ResolveServers(time.Second * 3)
	return len(others) > 0
}
