package main

import "github.com/draeron/gopkgs/logger"

type natsLogger struct {
	*logger.SugaredLogger
}

func (n *natsLogger) Noticef(format string, v ...interface{}) {
	n.SugaredLogger.Infof(format, v...)
}

func (n *natsLogger) Warnf(format string, v ...interface{}) {
	n.SugaredLogger.Warnf(format, v...)
}

func (n *natsLogger) Fatalf(format string, v ...interface{}) {
	n.SugaredLogger.Fatalf(format, v...)
}

func (n *natsLogger) Errorf(format string, v ...interface{}) {
	n.SugaredLogger.Errorf(format, v...)
}

func (n *natsLogger) Debugf(format string, v ...interface{}) {
	n.SugaredLogger.Debugf(format, v...)
}

func (n *natsLogger) Tracef(format string, v ...interface{}) {
	n.SugaredLogger.Debugf(format, v...)
}
