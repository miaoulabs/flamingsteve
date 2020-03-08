package main

import (
	"flamingsteve/cmd"
	"flamingsteve/pkg/display"
	"flamingsteve/pkg/muthur"
	"flamingsteve/pkg/pimoroni5x5"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
	"github.com/spf13/pflag"
)

const (
	Width      = 5
	Height     = Width * PanelCount
	PanelCount = 2
)

var args = struct {
	name  string
	model string
}{}

func init() {
	pflag.StringVarP(&args.name, "name", "n", "matrix", "name for this display")
	pflag.StringVarP(&args.model, "model", "p", "pimoroni-5x10", "name for this display")
}

func main() {
	pflag.Parse()
	cmd.SetupLoggers()
	muthur.Connect(args.name)
	defer muthur.Close()

	log := logger.New("main")

	bus := cmd.InitI2C()
	defer bus.Close()

	d1, err := pimoroni5x5.New(bus, pimoroni5x5.I2C_DEFAULT_ADDRESS)
	cmd.Must(err)
	d1.Clear(color.Black)
	defer d1.Clear(color.Green)

	d2, err := pimoroni5x5.New(bus, pimoroni5x5.I2C_ALTERNATE_ADDRESS)
	cmd.Must(err)
	d2.Clear(color.Black)
	defer d2.Clear(color.Green)

	remote := display.NewListener(args.name, args.model, Width, Height, func(msg *display.Message) {
		//log.Infof("received display instructions, length: %v", len(msg.Pixels))

		for i, p := range msg.Pixels {
			x := i % Width
			y := i / Width

			d := d2
			if y >= Width {
				d = d1
				y = y % Width
			}

			// hflip
			x = Width - 1 - x

			if p == '1' {
				d.SetPixel(y, x, color.Red)
			} else {
				d.SetPixel(y, x, color.Black)
			}
		}

		err = d1.Show()
		if err != nil {
			log.Errorf("failed to send to display #1: %v", err)
		}
		err = d2.Show()
		if err != nil {
			log.Errorf("failed to send to display #2: %v", err)
		}

	})
	defer remote.Close()

	cmd.WaitForCtrlC()
}
