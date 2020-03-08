package cmd

import (
	"flamingsteve/pkg/ak9753"
	ak9753presence "flamingsteve/pkg/ak9753/presence"
	"flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/amg8833"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/grpc"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/pimoroni5x5"
	"flamingsteve/pkg/presence"
	"fmt"
	"os"
	"os/signal"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
	"syscall"

	logger2 "flamingsteve/pkg/logger"
	"github.com/draeron/gopkgs/logger"
)

var (
	TerminationSignals = []os.Signal{
		syscall.SIGKILL, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT,
	}
)

func InitI2C() i2c.BusCloser {
	log := NewLogger("i2c")

	state, err := host.Init()
	Must(err)

	for i, drv := range state.Loaded {
		log.Infof("driver #%d: %v", i, drv.String())
	}

	bus, err := i2creg.Open("")
	Must(err)

	log.Infof("i2c bus %s is open", bus.String())

	return bus
}

func WaitForCtrlC() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, TerminationSignals...)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
}

func Must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func SetupLoggers() {
	ak9753.SetLoggerFactory(NewLogger)
	ak9753presence.SetLoggerFactory(NewLogger)
	amg8833.SetLoggerFactory(NewLogger)
	presence.SetLoggerFactory(NewLogger)
	remote.SetLoggerFactory(NewLogger)
	muthur.SetLoggerFactory(NewLogger)
	discovery.SetLoggerFactory(NewLogger)
	pimoroni5x5.SetLoggerFactory(NewLogger)
	grpc.SetLoggerFactory(NewLogger)
}

func NewLogger(name string) logger2.Logger {
	return logger.New(name)
}
