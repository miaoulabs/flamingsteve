package main

import (
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

var (
	threshold = pflag.Float32P("threshold", "t", 10, "presence threshold")
	interval  = pflag.DurationP("interval", "i", time.Millisecond*30, "interval for IR evaluration")
	smoothing = pflag.Float32P("smoothing", "s", 0.05, "0.3 very steep, 0.1 less steep, 0.05 less steep")
	ui        = pflag.Bool("ui", false, "display real time information on the terminal")
	orphan    = pflag.Bool("orphan", false, "don't try to connect to muthur")
	natsUrl   = pflag.String("nats-server", "", "publish nats server where to push the sensor data")
)

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

	ak9753.SetLogger(logger.New("ak9753"))

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

	if !*orphan {
		log.Infof("adoption mode enabled, scanning for mother")

		if *natsUrl == "" {
			var mothers discovery.Servers

			// keep trying until there is a connections
			for mothers == nil || len(mothers) == 0 {
				mothers = discovery.ResolveServers(time.Second * 2)
			}

			//log.Infof("found %d muthur on the local network", len(mothers))
			//*natsUrl = fmt.Sprintf("nats://%s:%d", mothers[0].HostName, mothers[0].Port)
			*natsUrl = mothers[0].HostName
		}

		log.Infof("connecting to muthur at '%s'", *natsUrl)
		device, err = remote.NewPublisher(device, *natsUrl)
		log.StopIfErr(err)
	}
	defer device.Close()

	detector = presence.New(device, &presence.Options{
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
