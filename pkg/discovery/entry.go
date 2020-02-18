package discovery

import (
	"flamingsteve/pkg/muthur"
	"net"
	"time"
)

type Entry struct {
	Type     EntryType `json:"type"` // sensor, display, detector
	Name     string    `json:"name"`
	IP       net.IP    `json:"ip"`
	Hostname string    `json:"hostname"`

	DataTopic   string `json:"raw_topic"` // raw data
	ConfigTopic string `json:"config_topic"`
	Model       string `json:"model"`
	Id          string `json:"id"`
}

const (
	Sensor = EntryType("sensor")
	Detector = EntryType("detector")
	Display = EntryType("display")
)

type EntryType string

func (e *Entry) FetchConfig() ([]byte, error) {
	resp, err := muthur.Connection().Request(e.ConfigTopic, []byte{}, time.Millisecond*500)
	if err != nil {
		return nil, err
	}
	return resp.Data, err
}

func (e *Entry) PushConfig(data []byte) error {
	return muthur.EncodedConnection().Publish(e.ConfigTopic, data)
}
