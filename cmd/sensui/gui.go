package main

import (
	"flamingsteve/pkg/ak9753/presence"
	"fmt"
	"image"
	"sync"
	"time"

	"flamingsteve/pkg/ak9753"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/draeron/gopkgs/logger"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type gui struct {
	smoothing float64
	presence  int
	base      float64
	ir        [ak9753.FieldCount][]float64
	irTime    [ak9753.FieldCount][]time.Time
	log       *logger.SugaredLogger
	changed   chan bool
	dev       ak9753.Device
	wnd       nucular.MasterWindow

	selectedSensorIndex int

	sync.RWMutex
}

const (
	LeftCenter = "LC"
	WidgetFlags = nucular.WindowNoScrollbar | nucular.WindowNoHScrollbar
	Height = 35
)

func (u *gui) updateSensorData() {
	maxValue := 60

	for range u.changed {
		if u.wnd.Closed() { //quit
			return
		}
		u.Lock()
		for i := 0; i < len(u.ir); i++ {
			if len(u.ir[i]) > maxValue {
				u.ir[i] = u.ir[i][1:]
				u.irTime[i] = u.irTime[i][1:]
			}
			ir, _ := u.dev.IR(i)
			u.ir[i] = append(u.ir[i], float64(ir))
			u.irTime[i] = append(u.irTime[i], time.Now())

			//if len(u.ir[i]) > maxValue {
			//	u.ir[i] = u.ir[i][:maxValue]
			//	u.irTime[i] = u.irTime[i][:maxValue]
			//}
			//u.ir[i] = append([]float64{float64(detector.IR(i))}, u.ir[i]...)
			//u.irTime[i] = append([]time.Time{time.Now()}, u.irTime[i]...)
		}
		u.Unlock()

		u.wnd.Changed() // force redraw
	}
}

func (ui *gui) renderUi(w *nucular.Window) {
	w.Row(SensorWidth).Static()

	w.LayoutSetWidth(SensorWidth)
	if s := w.GroupBegin("Sensors", WidgetFlags); s != nil {
		if ui.selectedSensorIndex > 0 {
			ui.renderSensors(s)
		}
		s.GroupEnd()
	}

	w.LayoutFitWidth(0, w.LayoutAvailableWidth()-SensorWidth)
	if p := w.GroupBegin("Properties", WidgetFlags); p != nil {
		ui.renderProps(p)
		p.GroupEnd()
	}
}

func (ui *gui) selectSensor(id string) {
	if sensor := sensors.Get(id); sensor != nil {
		ui.Lock()
		if ui.dev != nil {
			ui.dev.Unsubscribe(ui.changed)
		}

		if sensor != nil {
			sensor.Device.Subscribe(ui.changed)
			ui.dev = sensor.Device
		}
		ui.Unlock()
	}
}

func (ui *gui) renderSensors(w *nucular.Window) {
	w.RowScaled(w.LayoutAvailableHeight() / 3).Dynamic(3)

	w.Spacing(1)
	ui.drawSensor(w, 0)
	w.Spacing(1)
	ui.drawSensor(w, 1)
	w.Spacing(1)
	ui.drawSensor(w, 2)
	w.Spacing(1)
	ui.drawSensor(w, 3)
}

func (ui *gui) renderProps(w *nucular.Window) {
	//double := func() {
	//	w.Row(Height).Dynamic(2)
	//}
	single := func() {
		w.Row(Height).Dynamic(1)
	}

	sensorIds := append([]string{"none"}, sensors.Keys()...)
	single()

	oldSelected := ui.selectedSensorIndex
	ui.selectedSensorIndex = w.ComboSimple(sensorIds, ui.selectedSensorIndex, Height)
	if oldSelected != ui.selectedSensorIndex && ui.selectedSensorIndex > 0 {
		ui.selectSensor(sensorIds[ui.selectedSensorIndex])
	}

	if ui.selectedSensorIndex == 0 {
		return
	}

	name := sensorIds[ui.selectedSensorIndex]

	sensor := sensors.Get(name)
	if sensor == nil {
		return
	}

	padding := Height/2

	w.Spacing(1)
	w.Row(Height * 4 + padding).Dynamic(1)
	if p := w.GroupBegin("Properties", WidgetFlags | nucular.WindowBorder); p != nil {
		ui.renderAk9753SensorData(p, *sensor)
		p.GroupEnd()
	}

	w.Row(padding).Dynamic(1)
	w.Spacing(1)

	w.Row(Height * 6 + padding).Dynamic(1)
	if p := w.GroupBegin("Properties", WidgetFlags| nucular.WindowBorder); p != nil {
		ui.renderAk9753Detector(p, *sensor)
		p.GroupEnd()
	}
}

func (ui *gui) drawSensor(w *nucular.Window, idx int) {
	bounds, out := w.Custom(nstyle.WidgetStateActive)
	if out == nil {
		return
	}

	ui.RLock()
	irTime := ui.irTime[idx]
	ir := ui.ir[idx]
	ui.RUnlock()

	if len(ir) > 2 {
		ts := &chart.TimeSeries{
			XValues: irTime,
			YValues: ir,
			Style: chart.Style{
				Show:        true,
				StrokeColor: drawing.ColorBlack,
				StrokeWidth: 2,
			},
		}

		//maxSeries := &chart.MaxSeries{
		//	Style: chart.Style{
		//		Show:            true,
		//		StrokeColor:     chart.ColorAlternateGray,
		//		StrokeDashArray: []float64{5.0, 5.0},
		//	},
		//	Name: "max",
		//	InnerSeries: ts,
		//}

		bgcol := drawing.ColorBlue

		if ir[len(ir)-1] > ui.base {
			bgcol = drawing.ColorRed
		}

		//if detector.PresentInField(idx) {
		//	bgcol = drawing.ColorRed
		//}

		graph := &chart.Chart{
			Width:  bounds.W,
			Height: bounds.H,
			Title:  fmt.Sprintf("IR%d", idx),
			Series: []chart.Series{
				ts,
				//chart.LastValueAnnotation(maxSeries),
				chart.LastValueAnnotation(ts),
			},
			Background: chart.Style{
				Show:      true,
				FillColor: bgcol,
			},
			Canvas: chart.Style{
				Show:      true,
				FillColor: bgcol,
			},
			YAxis: chart.YAxis{
				Ascending: true,
				Zero: chart.GridLine{
					Value: 0,
					Style: chart.Style{
						Show:            true,
						StrokeColor:     chart.ColorAlternateGray,
						StrokeDashArray: []float64{5.0, 5.0},
					},
				},
				Range: &chart.ContinuousRange{
					Min: -400,
					Max: 1000,
				},
			},
		}

		collector := &chart.ImageWriter{}
		err := graph.Render(chart.PNG, collector)
		ui.log.LogIfErr(err)
		//gui.log.Errorf("error rendering graph: %v\n", err)

		img, err := collector.Image()
		if err == nil {
			if rgba, ok := img.(*image.RGBA); ok {
				out.DrawImage(bounds, rgba)
			}
		} else {
			ui.log.LogIfErr(err)
			//fmt.Fprintf(os.Stderr, "error collecting graph: %v\n", err)
		}
	}
}

func (ui *gui) renderAk9753SensorData(p *nucular.Window, sensor Sensor) {
	p.Row(Height).Dynamic(2)

	p.Label("Sensor Info", LeftCenter)
	p.Spacing(1)

	p.Label("Model: ", LeftCenter)
	p.Label(sensor.Ident.Model, LeftCenter)

	p.Label("Hostname: ", LeftCenter)
	p.Label(sensor.Ident.Hostname, LeftCenter)

	p.Label("IP: ", LeftCenter)
	p.Label(sensor.Ident.IP.String(), LeftCenter)
}

func (ui *gui) renderAk9753Detector(p *nucular.Window, sensor Sensor) {
	options := presence.UnmarshalOptions(sensor.Detector.Configs())
	updated := false


	p.Row(Height).Dynamic(1)
	p.Label("Presence Detector", LeftCenter)

	p.Row(Height).Dynamic(2)
	p.Label(fmt.Sprintf("Min Sensors: %d", options.MinimumSensors), LeftCenter)
	updated = p.SliderInt(1, &options.MinimumSensors, ak9753.FieldCount, 1)

	thresh := int(options.Threshold)
	p.Label(fmt.Sprintf("Threshold: %d", thresh), LeftCenter)
	updated = updated || p.SliderInt(-500, &thresh, 2000, 10)

	delay := int(options.Delay / time.Millisecond)
	p.Label(fmt.Sprintf("Delay: %d ms", delay), LeftCenter)
	updated = updated || p.SliderInt(10, &delay, 4000, 10)

	p.Label(fmt.Sprintf("Smoothing: %d", options.Smoothing), LeftCenter)
	updated = updated || p.SliderInt(2, &options.Smoothing, 20, 1)

	if updated {
		options.Threshold = float32(thresh)
		options.Delay = time.Duration(delay) * time.Millisecond
		sensor.Detector.SetConfigs(options.Marshal())

		if mean, ok := sensor.Device.(*ak9753.Mean); ok {
			mean.SetSampleCount(options.Smoothing)
		}
	}

	p.Row(Height).Dynamic(1)
	if p.ButtonText("Send Config") {
		// todo: send the new config to the sensor
	}
}
