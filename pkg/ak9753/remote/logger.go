package remote

import "flamingsteve/pkg/logger"

var log = logger.Dummy()

func SetLogger(newlogger logger.Logger) {
	log = newlogger
}