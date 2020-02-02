package main

import (
	"flamingsteve/pkg/ak9753"
	"fmt"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"image"
	"os"
	"sync"
	"time"
)

type ui struct {
	smoothing float64
	presence  int
	base      float64
	ir        [ak9753.FieldCount][]float64
	irTime    [ak9753.FieldCount][]time.Time

	sync.RWMutex
}

func (ui *ui) renderUi(w *nucular.Window) {
	w.Row(SensorWidth).Static()

	w.LayoutSetWidth(SensorWidth)
	if s := w.GroupBegin("Sensors", nucular.WindowNoScrollbar|nucular.WindowBorder); s != nil {
		ui.renderSensors(s)
		s.GroupEnd()
	}

	w.LayoutFitWidth(0, w.LayoutAvailableWidth()-SensorWidth)
	if p := w.GroupBegin("Properties", nucular.WindowDynamic|nucular.WindowBorder); p != nil {
		ui.renderProps(p)
		p.GroupEnd()
	}
}

func (ui *ui) renderSensors(w *nucular.Window) {
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

func (ui *ui) renderProps(p *nucular.Window) {
	height := 40
	double := func() {
		p.Row(height).Dynamic(2)
	}
	single := func() {
		p.Row(height).Dynamic(1)
	}

	opts := detector.Options()

	double()
	//r := 65536.0 / 2
	p.Label("Base line: ", "LT")
	p.Label(fmt.Sprintf("%f", ui.base), "RT")

	single()
	if p.SliderFloat(-1000, &ui.base, 2000, 3000/100) {
	}

	double()
	p.Label("Smoothing: ", "LT")
	p.Label(fmt.Sprintf("%f", ui.smoothing), "RT")

	single()
	if p.SliderFloat(0.01, &ui.smoothing, 0.5, 0.05) {
	}

	double()
	p.Label("Presence: ", "LT")
	p.Label(fmt.Sprintf("%d", ui.presence), "RT")

	single()
	if p.SliderInt(1, &ui.presence, 20, 1) {
	}

	opts.Smoothing = float32(ui.smoothing)
	opts.PresenceThreshold = float32(ui.presence)
	detector.SetOptions(opts)
}

func (ui *ui) drawSensor(w *nucular.Window, idx int) {
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
			Width: bounds.W,
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
						Show: true,
						StrokeColor:     chart.ColorAlternateGray,
						StrokeDashArray: []float64{5.0, 5.0},
					},
				},
				Range: &chart.ContinuousRange{
					Min: -300,
					Max: 8000,
				},
			},
		}

		collector := &chart.ImageWriter{}
		err := graph.Render(chart.PNG, collector)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error rendering graph: %v\n", err)
		}

		img, err := collector.Image()
		if err == nil {
			if rgba, ok := img.(*image.RGBA); ok {
				out.DrawImage(bounds, rgba)
			}
		} else {
			fmt.Fprintf(os.Stderr, "error collecting graph: %v\n", err)
		}
	}
}