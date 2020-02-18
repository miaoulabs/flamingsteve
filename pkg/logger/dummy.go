package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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
type LoggerFactory func(name string) Logger

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

func CurrentPackageName() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	pkage := ""
	funcName := parts[pl-1]
	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		pkage = strings.Join(parts[0:pl-2], ".")
	} else {
		pkage = strings.Join(parts[0:pl-1], ".")
	}
	pkage = strings.TrimPrefix(pkage, "flamingsteve/")
	pkage = strings.TrimPrefix(pkage, "pkg/")
	return pkage
}
