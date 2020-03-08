package main

import (
	"flamingsteve/cmd"
	"flamingsteve/pkg/ak9753"
	ak9753presence "flamingsteve/pkg/ak9753/presence"
	"flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/presence"

	"periph.io/x/periph/conn/i2c"
)

func start_ak9753(bus i2c.BusCloser) {
	ak, err := ak9753.New(bus, ak9753.I2C_DEFAULT_ADDRESS)
	log.StopIfErr(err)

	if ak == nil {
		log.Fatal("null device")
	}

	var detector *ak9753presence.Detector
	var device ak9753.Device

	device, err = ak9753.NewReader(ak)
	log.StopIfErr(err)

	did, _ := device.DeviceId()
	cid, _ := device.CompagnyCode()
	log.Infof("device id: 0x%x, compagny id: 0x%x", did, cid)

	if !args.orphan {
		log.Infof("adoption mode enabled, scanning for muthur")
		device, err = remote.NewPublisher(device, sensorId)
		log.StopIfErr(err)
	}
	defer device.Close()

	if !args.noPresence {
		log.Infof("creating presence detector")

		detector, err = ak9753presence.New(device, nil)
		log.StopIfErr(err)

		pub := presence.NewPublisher(detector, detectId)
		defer pub.Close()
	}

	if args.ui {
		d := ak9753_display{
			closed:   make(chan bool),
			device:   device,
			detector: detector,
		}
		go d.textDisplay()
		defer d.close()
	}

	cmd.WaitForCtrlC()
}
