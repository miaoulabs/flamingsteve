package presence

import "flamingsteve/pkg/logger"

var log = logger.Dummy()

func SetLoggerFactory(newLogger logger.LoggerFactory) {
	log = newLogger(logger.CurrentPackageName())
}
