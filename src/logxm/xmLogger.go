package logxm

import (
	"os"
	"sync"

	logrus "github.com/sirupsen/logrus"
)

// XmLogger should be a sinle instance
type XmLogger struct {
	logger      *logrus.Logger
	addHostName bool
	hostname    string
	mutex       sync.Mutex
}

// hold instance of xmLogger.
var logger *XmLogger

// xmNew creates instance of XmLogger. This logger should be the log
func xmNew(addHostName bool) *XmLogger {
	logger = &XmLogger{logger: logrus.New(), addHostName: addHostName}
	if addHostName {
		// always add host name to log.
		hostname, _ := os.Hostname()
		logger.hostname = hostname
	}
	return logger
}

func (x *XmLogger) getEntry() *logrus.Entry {
	if x.addHostName {
		return x.logger.WithField("host", x.hostname)
	}
	return logrus.NewEntry(x.logger)
}

// Debug level logger
func (x *XmLogger) Debug(args ...interface{}) {
	x.getEntry().Debug(args...)
}

// Info level logger
func (x *XmLogger) Info(args ...interface{}) {
	x.getEntry().Info(args...)
}

// Error level logger
func (x *XmLogger) Error(args ...interface{}) {
	x.getEntry().Error(args...)
}

// Warn level logger
func (x *XmLogger) Warn(args ...interface{}) {
	x.getEntry().Warn(args...)
}

// Fatal level logger
func (x *XmLogger) Fatal(args ...interface{}) {
	x.getEntry().Fatal(args...)
}
