package logger

import (
	"fmt"
	"os"
)

/*
  This match the signature of go.uber.org/zap sugarred logger
*/
type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type dummy struct{}

func Dummy() Logger {
	return &dummy{}
}

func (d dummy) Debugf(template string, args ...interface{}) {
}

func (d dummy) Infof(template string, args ...interface{}) {
}

func (d dummy) Warnf(template string, args ...interface{}) {
}

func (d dummy) Errorf(template string, args ...interface{}) {
}

func (d dummy) DPanicf(template string, args ...interface{}) {
	panic(fmt.Sprintf(template, args))
}

func (d dummy) Panicf(template string, args ...interface{}) {
	panic(fmt.Sprintf(template, args))
}

func (d dummy) Fatalf(template string, args ...interface{}) {
	os.Exit(1)
}
