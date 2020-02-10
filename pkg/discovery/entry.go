package discovery

import "net"

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
