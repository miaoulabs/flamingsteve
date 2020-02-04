package discovery

import (
	"context"
	"github.com/draeron/gopkgs/logger"
	"github.com/grandcat/zeroconf"
	"strings"
	"time"
)

type Servers []*zeroconf.ServiceEntry

func ResolveServers(wait time.Duration) []*zeroconf.ServiceEntry {
	log := logger.New("zeroconf")

	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.StopIfErr(err)
	}

	dnses := []*zeroconf.ServiceEntry{}

	entries := make(chan *zeroconf.ServiceEntry, 10)

	//go func(results <-chan *zeroconf.ServiceEntry) {
	//
	//	log.Info("no more entries")
	//}(entries)

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
