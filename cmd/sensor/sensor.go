package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"flamingsteve/pkg/ak9753"
	rm "flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/presence_detector"
	"github.com/spf13/pflag"
	"periph.io/x/periph"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"

	"time"
)

var (
	threshold = pflag.Float32P("threshold", "t", 10, "presence threshold")
	interval  = pflag.DurationP("interval", "i", time.Millisecond*30, "interval for IR evaluration")
	smoothing = pflag.Float32P("smoothing", "s", 0.05, "0.3 very steep, 0.1 less steep, 0.05 less steep")
	ui        = pflag.Bool("ui", false, "display real time informatio on the terminal")
	publish   = pflag.BoolP("publish", "p", false, "url for publish data push")
	remote    = pflag.Bool("remote", false, "connect to a remote sensor")
	natsUrl   = pflag.String("nats-server", "", "publish nats server where to push the sensor data")
)

func hostInit() (*periph.State, error) {
	return host.Init()
}

var detector *pdetect.Detector

func mainImpl() error {

	var err error
	var device ak9753.Device

	if !*remote {
		state, err := hostInit()
		if err != nil {
			return err
		}

		for i, drv := range state.Loaded {
			fmt.Printf("driver #%d: %v\n", i, drv.String())
		}

		b, err := i2creg.Open("")
		if err != nil {
			log.Fatal(err)
		}
		defer b.Close()

		fmt.Printf("i2c bus %s is open\n", b.String())

		ak, err := ak9753.New(b, ak9753.I2C_DEFAULT_ADDRESS)
		if err != err {
			return err
		}

		if ak == nil {
			return errors.New("null device")
		}

		device, err = ak9753.NewReader(ak)
		if err != nil {
			return err
		}
	} else {
		device, err = rm.NewSuscriber(*natsUrl)
		if err != nil {
			return err
		}
	}

	did, _ := device.DeviceId()
	cid, _ := device.CompagnyCode()
	fmt.Printf("device id: 0x%x, compagny id: 0x%x\n", did, cid)

	if *publish {
		device, err = rm.NewPublisher(device, *natsUrl)
		if err != nil {
			return err
		}
	}
	defer device.Close()

	detector = pdetect.New(device, &pdetect.Options{
		Interval:          *interval,
		PresenceThreshold: *threshold,
		MovementThreshold: 10,
		Smoothing:         *smoothing,
	})
	defer detector.Close()

	if *ui {
		d := display{closed: make(chan bool)}
		go d.textDisplay()
		defer d.close()
	}

	waitForTerm()


	return err
}

func main() {
	pflag.Parse()

	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "i2c-io: %s.\n", err)
		os.Exit(1)
	}
}

func waitForTerm() {
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
