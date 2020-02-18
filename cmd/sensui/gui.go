package main

import (
	"fmt"
	"image"
	"sync"
	"time"

	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/ak9753/presence"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/draeron/gopkgs/logger"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type gui struct {
	ir        [ak9753.FieldCount][]float64
	irTime    [ak9753.FieldCount][]time.Time
	log       *logger.SugaredLogger
	changed   chan bool
	wnd       nucular.MasterWindow

	selectedSensorIndex int
	currentSensor *Sensor

	options presence.Options

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
			ir, _ := u.currentSensor.Device.IR(i)
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
		if ui.currentSensor != nil {
			ui.currentSensor.Device.Unsubscribe(ui.changed)
		}

		sensor.Device.Subscribe(ui.changed)
		ui.currentSensor = sensor
		ui.options = presence.UnmarshalOptions(sensor.LocalDetector.Configs())
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

	sensorIds := append([]string{"none"}, sensors.Keys()...)
	w.Row(Height).Dynamic(1)

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
		ui.renderAk9753SensorData(p)
		p.GroupEnd()
	}

	w.Row(padding).Dynamic(1)
	w.Spacing(1)

	w.Row(Height * 6 + padding).Dynamic(1)
	if p := w.GroupBegin("Properties", WidgetFlags| nucular.WindowBorder); p != nil {
		ui.renderAk9753Detector(p)
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

		if ir[len(ir)-1] > float64(ui.options.Threshold) {
			bgcol = drawing.ColorRed
		}

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
				Style: chart.Style{
					Show: true,
					StrokeColor: chart.ColorTransparent,
					FontColor: chart.ColorTransparent,
				},
				GridMajorStyle: chart.Style{
					Show:            true,
					StrokeColor:     chart.ColorWhite,
					StrokeDashArray: []float64{5.0, 5.0},
					StrokeWidth: 1,
				},
				GridLines: []chart.GridLine{
					{
						Value:   float64(ui.options.Threshold),
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

func (ui *gui) renderAk9753SensorData(p *nucular.Window) {
	p.Row(Height).Dynamic(2)

	p.Label("Sensor Info", LeftCenter)
	p.Spacing(1)

	p.Label("Model: ", LeftCenter)
	p.Label(ui.currentSensor.Ident.Model, LeftCenter)

	p.Label("Hostname: ", LeftCenter)
	p.Label(ui.currentSensor.Ident.Hostname, LeftCenter)

	p.Label("IP: ", LeftCenter)
	p.Label(ui.currentSensor.Ident.IP.String(), LeftCenter)
}

func (ui *gui) renderAk9753Detector(p *nucular.Window) {
	updated := false

	p.Row(Height).Dynamic(1)
	p.Label("Presence LocalDetector", LeftCenter)

	p.Row(Height).Dynamic(2)
	p.Label(fmt.Sprintf("Min Sensors: %d", ui.options.MinimumSensors), LeftCenter)
	updated = p.SliderInt(1, &ui.options.MinimumSensors, ak9753.FieldCount, 1)

	thresh := int(ui.options.Threshold)
	p.Label(fmt.Sprintf("Threshold: %d", thresh), LeftCenter)
	updated = updated || p.SliderInt(-500, &thresh, 2000, 5)

	delay := int(ui.options.Delay / time.Millisecond)
	p.Label(fmt.Sprintf("Delay: %d ms", delay), LeftCenter)
	updated = updated || p.SliderInt(10, &delay, 4000, 5)

	p.Label(fmt.Sprintf("Smoothing: %d", ui.options.Smoothing), LeftCenter)
	updated = updated || p.SliderInt(2, &ui.options.Smoothing, 20, 1)

	if updated {
		ui.options.Threshold = float32(thresh)
		ui.options.Delay = time.Duration(delay) * time.Millisecond
		ui.currentSensor.LocalDetector.SetConfigs(ui.options.Marshal())

		if mean, ok := ui.currentSensor.Device.(*ak9753.Mean); ok {
			mean.SetSampleCount(ui.options.Smoothing)
		}
	}

	p.Row(Height).Dynamic(1)
	if p.ButtonText("Send Config") {
		ui.currentSensor.RemoteDetector.SetConfigs(ui.options.Marshal())
	}
}
