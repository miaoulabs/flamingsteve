package main

import (
	"flamingsteve/cmd"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/draeron/gopkgs/logger"
	"github.com/spf13/pflag"
)

var args = struct {
	threshold  float32
	duration   time.Duration
	noPresence bool
	ui         bool
	orphan     bool
	mean       int
	name       string

	stype SensorType
}{
	stype: SensorTypeNone,
}

func init() {

	pflag.Float32VarP(&args.threshold, "threshold", "t", 100, "presence threshold")
	pflag.IntVar(&args.mean, "mean", 6, "number of sample to use for mean")
	pflag.DurationVar(&args.duration, "duration", time.Millisecond*100, "time before a presence is considered")
	pflag.BoolVar(&args.noPresence, "no-presence", false, "disable presence detector")

	pflag.BoolVar(&args.ui, "ui", false, "display real time information on the terminal")
	pflag.BoolVar(&args.orphan, "orphan", false, "don't try to connect to muthur")

	pflag.StringVarP(&args.name, "name", "n", "", "sensor name used for discovery")

	pflag.Var(&args.stype, "type", fmt.Sprintf("sensor model [%s]", strings.Join(SensorTypeNames()[1:], ", ")))
}

var (
	log      *logger.SugaredLogger
	sensorId *discovery.Component
	detectId *discovery.Component
)

func main() {
	pflag.Parse()
	cmd.SetupLoggers()

	if args.stype < 0 {
		fmt.Fprintf(os.Stderr, "no sensor type specified, can be [%s]\n", strings.Join(SensorTypeNames()[1:], ", "))
		pflag.PrintDefaults()
		os.Exit(1)
	}

	log = logger.New("main")

	log.Infof("sensor started")
	defer log.Infof("sensor stopped")

	bus := cmd.InitI2C()
	defer bus.Close()

	if args.name == "" {
		args.name = args.stype.String()
	}

	if !args.orphan {
		sensorId = discovery.NewComponent(discovery.IdentifierConfig{
			Name:  args.name,
			Model: args.stype.String(),
			Type:  discovery.Sensor,
		})

		muthur.Connect(sensorId.Id())
		sensorId.Connect()

		defer sensorId.Close()
		defer muthur.Close()

		if !args.noPresence {
			detectId = discovery.NewComponent(discovery.IdentifierConfig{
				Name:  args.name,
				Model: args.stype.String(),
				Type:  discovery.Detector,
			})
		}
	}

	switch args.stype {
	case SensorTypeAk9753:
		start_ak9753(bus)
	case SensorTypeAmg8833:
		start_amg8833(bus)
	}
}
