package main

import (
	"flamingsteve/pkg/ak9753"
	ak9753presence "flamingsteve/pkg/ak9753/presence"
	"fmt"
	"time"

	tm "github.com/buger/goterm"
)

type ak9753_display struct {
	closed   chan bool
	detector *ak9753presence.Detector
	device   ak9753.Device
}

func center(text string, width int) string {
	return fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(text))/2, text))
}

func (d *ak9753_display) close() {
	d.closed <- true
	close(d.closed)
}

func (d *ak9753_display) textDisplay() {
	width := 8
	toXO := func(v bool) string {
		if v {
			return center("YES", width)
		} else {
			return center("no", width)
		}
	}

	time.Sleep(time.Millisecond * 1000)

	tm.Clear()

	tick := time.NewTicker(time.Millisecond * 100)
	defer tick.Stop()

	start := time.Now()

	for range tick.C {
		select {
		case <-d.closed:
			return // exit loop
		default:
		}

		tm.MoveCursor(1, 1)

		tm.Printf("            | %s | %s | %s | %s |\n",
			center("IR1", width),
			center("IR2", width),
			center("IR3", width),
			center("IR4", width),
		)

		if d.detector != nil {
			tm.Printf("presence    | %s | %s | %s | %s |\n",
				toXO(d.detector.PresentInField1()),
				toXO(d.detector.PresentInField2()),
				toXO(d.detector.PresentInField3()),
				toXO(d.detector.PresentInField4()),
			)
		}

		ir1, _ := d.device.IR1()
		ir2, _ := d.device.IR2()
		ir3, _ := d.device.IR3()
		ir4, _ := d.device.IR4()
		tm.Printf("sensor      | %8.3f | %8.3f | %8.3f | %8.3f |\n",
			ir1,
			ir2,
			ir3,
			ir4,
		)

		tmp, _ := d.device.Temperature()
		tm.Printf("temperature | %8.2f C\n", tmp)
		tm.Printf("elapsed     | %v\n", time.Now().Sub(start))
		tm.Flush()
	}
}
