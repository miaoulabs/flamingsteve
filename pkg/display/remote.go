package display

import (
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
)

type Remote struct {
	 ident discovery.Entry
}

func NewRemote(ident discovery.Entry) *Remote {
	r := &Remote{
		ident: ident,
	}
	return r
}

func (r *Remote) Draw(msg Message) error {
	return muthur.EncodedConnection().Publish(r.ident.DataTopic, msg)
}
