package main

import (
	"time"

	"flamingsteve/pkg/amg8833"
	tm "github.com/buger/goterm"
)

type amg8833_display struct {
	closed chan bool
	device amg8833.Device
}

func (d *amg8833_display) close() {
	d.closed <- true
	close(d.closed)
}

func (d *amg8833_display) textDisplay() {
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

		temps := d.device.Temperatures()

		for i := 0; i < amg8833.ROW_COUNT; i++ {
			tm.Printf("| %8.3f | %8.3f | %8.3f | %8.3f | %8.3f | %8.3f | %8.3f | %8.3f |\n",
				temps[i*amg8833.ROW_COUNT],
				temps[i*amg8833.ROW_COUNT+1],
				temps[i*amg8833.ROW_COUNT+2],
				temps[i*amg8833.ROW_COUNT+3],
				temps[i*amg8833.ROW_COUNT+4],
				temps[i*amg8833.ROW_COUNT+5],
				temps[i*amg8833.ROW_COUNT+6],
				temps[i*amg8833.ROW_COUNT+7],
			)
		}

		tm.Printf("elapsed: %v\n", time.Now().Sub(start))
		tm.Flush()
	}
}
