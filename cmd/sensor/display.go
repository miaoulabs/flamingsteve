package main

import (
	"fmt"
	tm "github.com/buger/goterm"
	"time"
)

type display struct {
	closed chan bool
}

func center(text string, width int) string {
	return fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(text))/2, text))
}

func (d *display) close() {
	d.closed <- true
	close(d.closed)
}

func (d *display) textDisplay() {
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
		tm.Printf("presence    | %s | %s | %s | %s |\n",
			toXO(detector.PresentInField1()),
			toXO(detector.PresentInField2()),
			toXO(detector.PresentInField3()),
			toXO(detector.PresentInField4()),
		)
		tm.Printf("sensor      | %8.2f | %8.2f | %8.2f | %8.2f |\n",
			detector.IR1(),
			detector.IR2(),
			detector.IR3(),
			detector.IR4(),
		)
		tm.Printf("derivative  | %8.2f | %8.2f | %8.2f | %8.2f |\n",
			detector.DerivativeOfIR1(),
			detector.DerivativeOfIR2(),
			detector.DerivativeOfIR3(),
			detector.DerivativeOfIR4(),
		)
		tm.Printf("temperature | %8.2f C\n", detector.Temperature())
		tm.Printf("elapsed     | %v\n", time.Now().Sub(start))
		tm.Flush()
	}
}
