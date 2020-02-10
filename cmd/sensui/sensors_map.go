package main

import (
	"flamingsteve/pkg/presence"
	"sort"
	"sync"

	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/discovery"
)

type SensorsMap struct {
	m sync.Map
}

type Sensor struct {
	Ident    discovery.Entry
	Device   ak9753.Device
	Detector presence.Detector
}

func (s *SensorsMap) Get(id string) *Sensor {
	if it, ok := s.m.Load(id); ok {
		if dev, ok := it.(*Sensor); ok {
			return dev
		}
	}
	return nil
}

func (s *SensorsMap) Set(id string, sensor Sensor) {
	s.m.Store(id, &sensor)
}

func (s *SensorsMap) Delete(id string) {
	s.m.Delete(id)
}

func (s *SensorsMap) Keys() []string {
	ids := []string{}

	s.m.Range(func(key, value interface{}) bool {
		ids = append(ids, key.(string))
		return true
	})

	sort.Strings(ids)
	return ids
}
