package main

import (
	"flamingsteve/pkg/muthur"
	"fmt"
	"os"
	"os/signal"
	"time"

	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/ak9753/presence"
	"flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/discovery"

	"github.com/draeron/gopkgs/logger"
	"github.com/spf13/pflag"
	"periph.io/x/periph"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

var args = struct {
	threshold float32
	interval  time.Duration
	presence  bool
	ui        bool
	orphan    bool
	mean      int
}{}

func init() {
	pflag.Float32VarP(&args.threshold, "threshold", "t", 100, "presence threshold")
	pflag.IntVar(&args.mean, "mean", 6, "number of sample to use for mean")
	pflag.DurationVarP(&args.interval, "interval", "i", time.Millisecond*30, "interval for presence evaluation")
	pflag.BoolVarP(&args.presence, "presence", "p", true, "test for presence")

	pflag.BoolVar(&args.ui, "ui", false, "display real time information on the terminal")
	pflag.BoolVar(&args.orphan, "orphan", false, "don't try to connect to muthur")

}

func hostInit() (*periph.State, error) {
	return host.Init()
}

var detector *presence.Detector

func main() {
	pflag.Parse()
	log := logger.New("main")
	presence.SetLogger(logger.New("detector"))
	remote.SetLogger(logger.New("remote"))
	ak9753.SetLogger(logger.New("ak9753"))
	muthur.SetLogger(logger.New("muthur"))

	var err error
	var device ak9753.Device

	state, err := hostInit()
	log.StopIfErr(err)

	for i, drv := range state.Loaded {
		log.Infof("driver #%d: %v", i, drv.String())
	}

	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	log.Infof("i2c bus %s is open", b.String())

	ak, err := ak9753.New(b, ak9753.I2C_DEFAULT_ADDRESS)
	log.StopIfErr(err)

	if ak == nil {
		log.Fatal("null device")
	}

	device, err = ak9753.NewReader(ak)
	log.StopIfErr(err)

	did, _ := device.DeviceId()
	cid, _ := device.CompagnyCode()
	log.Infof("device id: 0x%x, compagny id: 0x%x", did, cid)

	if !args.orphan {
		log.Infof("adoption mode enabled, scanning for muthur")

		ident := discovery.NewIdentifier(discovery.IdentifierConfig{
			Name:  "protopi",
			Model: "ak9753",
			Type:  "sensor",
		}, logger.New("identifier"))

		muthur.Connect(ident.Id())
		ident.Connect()

		defer ident.Close()
		defer muthur.Close()

		device, err = remote.NewPublisher(device, ident)
		log.StopIfErr(err)
	}
	defer device.Close()

	if args.presence {
		//detector = presence.New(device, &presence.Options{
		//	Interval:          args.interval,
		//	PresenceThreshold: args.threshold,
		//	Smoothing:         args.mean,
		//})
		//defer detector.Close()
	}

	if args.ui {
		d := display{closed: make(chan bool)}
		go d.textDisplay()
		defer d.close()
	}

	waitForTerm()
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
