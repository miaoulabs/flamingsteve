package display

import (
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/muthur"
	jsoniter "github.com/json-iterator/go"
	"sync"
)

type Remote struct {
	ident discovery.Entry
	info  Info
	mutex sync.RWMutex
}

func NewRemote(ident discovery.Entry) (*Remote, error) {
	r := &Remote{
		ident: ident,
	}

	err := r.UpdateInfo()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Remote) Draw(msg Message) error {
	return muthur.EncodedConnection().Publish(r.ident.DataTopic, msg)
}

func (r *Remote) UpdateInfo() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	data, err := r.ident.FetchConfig()
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(data, &r.info)
	if err != nil {
		return err
	}

	return nil
}

func (r *Remote) Dimension() (w, h uint) {
	r.mutex.RLock()
	r.mutex.RUnlock()
	return r.info.Width, r.info.Height
}

func (r *Remote) Ident() discovery.Entry {
	return r.ident
}
