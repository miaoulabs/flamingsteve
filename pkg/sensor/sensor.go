package sensor

import "flamingsteve/pkg/discovery"

type Sensor interface {
	Ident() discovery.Entry
	Raw() []byte
	Unstructured() map[string]interface{}
}

type Remote struct {
	ident discovery.Entry
	data  []byte
}

func NewRemote(ident discovery.Entry) {

}
