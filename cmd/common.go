package cmd

import (
	"flamingsteve/pkg/ak9753"
	ak9753presence "flamingsteve/pkg/ak9753/presence"
	"flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/presence"
	"fmt"
	"os"
	"os/signal"

	logger2 "flamingsteve/pkg/logger"
	"github.com/draeron/gopkgs/logger"
)

func WaitForCtrlC() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done

	println("stopping application")
}

func Must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func SetupLoggers() {
	ak9753.SetLoggerFactory(NewLogger)
	ak9753presence.SetLoggerFactory(NewLogger)
	presence.SetLoggerFactory(NewLogger)
	remote.SetLoggerFactory(NewLogger)
	muthur.SetLoggerFactory(NewLogger)
	discovery.SetLoggerFactory(NewLogger)
}

func NewLogger(name string) logger2.Logger {
	return logger.New(name)
}
