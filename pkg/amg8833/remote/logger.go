package remote

import (
	"flamingsteve/pkg/logger"
)

var log = logger.Dummy()

func SetLoggerFactory(newLogger func(name string) logger.Logger) {
	log = newLogger(logger.CurrentPackageName())
}
