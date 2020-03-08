package main

import (
	"flamingsteve/cmd"
	"flamingsteve/pkg/amg8833"
	"flamingsteve/pkg/amg8833/remote"
	"periph.io/x/periph/conn/i2c"
)

func start_amg8833(bus i2c.BusCloser) {
	phy, err := amg8833.New(bus, amg8833.I2C_DEFAULT_ADDRESS)
	log.StopIfErr(err)

	var device amg8833.Device

	device, err = amg8833.NewReader(phy)
	log.StopIfErr(err)

	if !args.orphan {
		log.Infof("adoption mode enabled, scanning for muthur")
		device, err = remote.NewPublisher(device, sensorId)
		log.StopIfErr(err)
	}
	defer device.Close()

	if args.ui {
		d := amg8833_display{
			closed: make(chan bool),
			device: device,
		}
		go d.textDisplay()
		defer d.close()
	}

	cmd.WaitForCtrlC()
}
