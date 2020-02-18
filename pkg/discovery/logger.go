package discovery

import "flamingsteve/pkg/logger"

var logFactory logger.LoggerFactory = func(name string) logger.Logger {
	return logger.Dummy()
}

func SetLoggerFactory(newLogger func(name string) logger.Logger) {
	logFactory = newLogger
}
