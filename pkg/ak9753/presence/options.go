package presence

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type Options struct {
	Delay          time.Duration `json:"delay"`           // duration before the detector is considered present
	Smoothing      int           `json:"smoothing"`       // number of data point to use for mean computing
	MinimumSensors int           `json:"minimum_sensors"` // number of sensor that need to
	Threshold      float32       `json:"threshold"`       // floor reading
}

const (
	DefaultDelay     = time.Millisecond * 100
	DefaultSmoothing = 5
	DefaultThreshold = float32(0)
	DefaultMinSensor = 2
)

func UnmarshalOptions(bytes []byte) Options {
	var opts Options

	if err := jsoniter.Unmarshal(bytes, &opts); err != nil {
		// todo: log err
	}

	return opts
}

func (o Options) Marshal() []byte {
	bytes, _ := jsoniter.Marshal(&o)
	return bytes
}

func (o Options) String() string {
	return fmt.Sprintf("Delay: %v, Smoothing: %v, Sensor Count: %v, Threshold: %v",
		o.Delay, o.Smoothing, o.MinimumSensors, o.Threshold,
	)
}
