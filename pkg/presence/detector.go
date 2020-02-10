package presence

import (
	"flamingsteve/pkg/notification"
)

type Detector interface {
	notification.Notifier

	// Return the configuration as a byte slice (should be json)
	Configs() []byte

	// Write configuration to the detector
	SetConfigs(data []byte)

	IsPresent() bool
	Close()
}

type DetectorState struct {
	Present bool `json:"present"`
}
