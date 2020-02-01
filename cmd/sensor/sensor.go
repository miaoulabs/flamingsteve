package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"fmt"

	tm "github.com/buger/goterm"

	"flamingsteve/pkg/ak9753"
	"flamingsteve/pkg/presence_detector"
	"periph.io/x/periph"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"

	"time"
)

func hostInit() (*periph.State, error) {
	return host.Init()
}

func mainImpl() error {

	state, err := hostInit()
	if err != nil {
		return err
	}

	for i, drv := range state.Loaded {
		fmt.Printf("driver #%d: %v\n", i, drv.String())
	}

	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	fmt.Printf("i2c bus %s is open\n", b.String())

	ak, err := ak9753.New(b, ak9753.I2C_DEFAULT_ADDRESS)
	if err != err {
		return err
	}

	if ak == nil {
		return errors.New("null device")
	}

	go func() {
		detector := pdetect.New(ak, &pdetect.Options{
			Interval:          time.Millisecond * 30,
			PresenceThreshold: 6,
			MovementThreshold: 10,
		})
		defer detector.Close()

		tick := time.NewTicker(time.Millisecond * 50)
		defer tick.Stop()

		width := 8

		toXO := func(v bool) string {
			if v {
				return center("YES", width)
			} else {
				return center("no", width)
			}
		}

		tm.Clear()

		for range tick.C {
			tm.MoveCursor(1,1)

			tm.Printf("             | %s | %s | %s | %s |\n",
				center("IR1", width),
				center("IR2", width),
				center("IR3", width),
				center("IR4", width),
			)
			tm.Printf("presence   : | %s | %s | %s | %s |\n",
				toXO(detector.PresentInField1()),
				toXO(detector.PresentInField2()),
				toXO(detector.PresentInField3()),
				toXO(detector.PresentInField4()),
			)
			tm.Printf("sensor     : | %8.2f | %8.2f | %8.2f | %8.2f |\n",
				detector.IR1(),
				detector.IR2(),
				detector.IR3(),
				detector.IR4(),
			)
			tm.Printf("derivative : | %8.2f | %8.2f | %8.2f | %8.2f |\n",
				detector.DerivativeOfIR1(),
				detector.DerivativeOfIR2(),
				detector.DerivativeOfIR3(),
				detector.DerivativeOfIR4(),
			)
			tmp, err := ak.Temperature()
			if err == nil {
				tm.Printf("temperature: | %8.2f C\n", tmp)
			}
			tm.Flush()
		}
	}()

	waitForTerm()

	return err
}

func center(text string, width int) string {
	return fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width + len(text))/2, text))
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "i2c-io: %s.\n", err)
		os.Exit(1)
	}
}

func waitForTerm() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done

	println("stopping application")

	os.Exit(0)
}
