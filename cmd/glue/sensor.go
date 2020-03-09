package main

import (
	"flamingsteve/pkg/ak9753"
	akremote "flamingsteve/pkg/ak9753/remote"
	"flamingsteve/pkg/amg8833"
	amgremote "flamingsteve/pkg/amg8833/remote"
	"flamingsteve/pkg/discovery"
)

type RawGetter interface {
	Raw() interface{}
}

type Sensor struct {
	RawGetter
	Ident discovery.Entry
}

func NewSensor(id discovery.Entry) *Sensor {
	switch id.Model {
	case ak9753.ModelName:
		dev, _ := akremote.NewSuscriber(id)
		return &Sensor{
			RawGetter: dev,
			Ident:     id,
		}

	case amg8833.ModelName:
		dev, _ := amgremote.NewSuscriber(id)
		return &Sensor{
			RawGetter: dev,
			Ident:     id,
		}

	}
	return nil
}
