package main

import (
	"github.com/sirupsen/logrus"
	"github.com/scaleway/taskor"
	"github.com/scaleway/taskor/runner/goroutine"
)

type logrusTaskor struct{}

var logrusLogger = logrus.WithField("WhoamI", "LogrusLogger")

// Debug -
func (l *logrusTaskor) Debug(msg string, extraFields map[string]interface{}) {
	logrusLogger.WithFields(extraFields).Debug(msg)
}

// Info -
func (l *logrusTaskor) Info(msg string, extraFields map[string]interface{}) {
	logrusLogger.WithFields(extraFields).Info(msg)
}

// Warn -
func (l *logrusTaskor) Warn(msg string, extraFields map[string]interface{}) {
	logrusLogger.WithFields(extraFields).Warn(msg)
}

// Error -
func (l *logrusTaskor) Error(msg string, extraFields map[string]interface{}) {
	logrusLogger.WithFields(extraFields).Error(msg)
}

func main() {
	taskor.SetLogger(&logrusTaskor{})

	config := goroutine.RunnerConfig{
		MaxBufferedMessage: 0,
	}
	taskor.New(goroutine.New(config))

}
