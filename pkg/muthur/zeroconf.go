package muthur

import (
	"context"
	"github.com/draeron/gopkgs/logger"
	"github.com/grandcat/zeroconf"
	"net"
	"strings"
	"time"
)

type Servers []*zeroconf.ServiceEntry

/*
	Will block until a muthur is found
*/
func MustResolveServer() *zeroconf.ServiceEntry {
	var svrs Servers

	wait := time.Millisecond * 1500

	for svrs == nil || len(svrs) == 0 {
		svrs = ResolveServers(wait)

		wait *= 2

		if wait > time.Second*8 {
			wait = time.Second * 8
		}
	}

	return svrs[0]
}

func ResolveServers(wait time.Duration) Servers {
	log := logger.New("zeroconf")

	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(zeroconf.SelectIfaces(ListMulticastInterfaces()))
	if err != nil {
		log.StopIfErr(err)
	}

	dnses := []*zeroconf.ServiceEntry{}

	entries := make(chan *zeroconf.ServiceEntry, 10)

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	log.Info("searching for muthur...")
	err = resolver.Browse(ctx, ZeroConfServiceName, ZeroConfDomain, entries)
	log.StopIfErr(err)

	<-ctx.Done()

	for entry := range entries {
		entry.HostName = strings.TrimSuffix(entry.HostName, ".local.")
		dnses = append(dnses, entry)
	}

	log.Infof("%d muthur server were found on the local network", len(dnses))

	return dnses
}

/*
	This a copy/paste from zeroconf package, except we filter
	vEthernet (docker) interface
*/
func ListMulticastInterfaces() []net.Interface {
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
