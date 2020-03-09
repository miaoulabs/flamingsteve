package main

import (
	"context"
	"flamingsteve/cmd"
	"flamingsteve/pkg/discovery"
	"flamingsteve/pkg/display"
	"flamingsteve/pkg/logger"
	"flamingsteve/pkg/muthur"
	"sync"
	"time"

	"flamingsteve/cmd/glue/api"
)

type Service struct {
	fpixels.FlamePixelsServer

	log logger.Logger

	displays map[string]*display.Remote
	sensors  map[string]Sensor

	mutex          sync.RWMutex
	displayScanner *discovery.Scanner
	sensorScanner  *discovery.Scanner
}

type Display struct {
	fpixels.Display
}

func NewService() fpixels.FlamePixelsServer {
	s := &Service{
		log:      cmd.NewLogger("fpixels"),
		displays: map[string]*display.Remote{},
		sensors: map[string]Sensor{},
	}

	muthur.Connect("glue")

	s.displayScanner = discovery.NewScanner(discovery.Display, s.onNewDisplay, s.onRmDisplay)
	s.displayScanner.Scan()                      // start a scan right away
	s.displayScanner.StartScan(time.Second * 30) // scan every 30 sec

	s.sensorScanner = discovery.NewScanner(discovery.Sensor, s.onNewSensor, s.onRmSensor)
	s.sensorScanner.Scan()
	s.sensorScanner.StartScan(time.Second * 30)

	return s
}

func (s *Service) Close() {
	s.displayScanner.Close()
}

func (s *Service) ListSensors(context.Context, *fpixels.EmptyRequest) (*fpixels.ListSensorsResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	resp := &fpixels.ListSensorsResponse{}

	for _, sen := range s.sensors {
		resp.Sensors = append(resp.Sensors, toDevice(sen.Ident))
	}

	return resp, nil
}

func (s *Service) GetSensorRawData(context.Context, *fpixels.SensorRawDataRequest) (*fpixels.SensorRawDataResponse, error) {
	panic("implement me")
}

func (s *Service) ListDisplays(ctx context.Context, req *fpixels.EmptyRequest) (*fpixels.ListDisplaysResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	resp := &fpixels.ListDisplaysResponse{}

	for _, v := range s.displays {
		w, h := v.Dimension()
		e := v.Ident()

		resp.Displays = append(resp.Displays, &fpixels.Display{
			Device: toDevice(e),
			Width:  uint32(w),
			Height: uint32(h),
		})
	}

	return resp, nil
}

func (s *Service) Draw(context.Context, *fpixels.DrawRequest) (*fpixels.DrawRequest, error) {
	panic("implement me")
}

//
//func (s *Service) Draw(ctx context.Context, req *fpixels.DrawRequest) (*fpixels.EmptyReply, error) {
//	disp := s.getDisplay(req.Id)
//	if disp == nil {
//		return nil, status.Errorf(codes.InvalidArgument, "display %v doesn't exists", req.Id)
//	}
//
//	//disp.Draw()
//
//	return &fpixels.EmptyReply{}, nil
//}

func (s *Service) getDisplay(id string) *display.Remote {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	disp, _ := s.displays[id]
	return disp
}

func (s *Service) onNewDisplay(entry discovery.Entry) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	remote, err := display.NewRemote(entry)
	if err != nil {
		s.log.Errorf("failed to connect to display: %v", err)
	}

	s.displays[entry.Id] = remote
}

func (s *Service) onRmDisplay(entry discovery.Entry) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.displays, entry.Id)
}

func (s *Service) onNewSensor(entry discovery.Entry) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sen := NewSensor(entry)
	if sen != nil {
		s.sensors[entry.Id] = *sen
	}
}

func (s *Service) onRmSensor(entry discovery.Entry) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sensors, entry.Id)
}

func toDevice(entry discovery.Entry) *fpixels.Device {
	return &fpixels.Device{
		Id:       entry.Id,
		Name:     entry.Name,
		Model:    entry.Model,
		Hostname: entry.Hostname,
	}
}
