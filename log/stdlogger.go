package log

import (
	"fmt"
	"log"
	"os"
)

// DefaultLogger - Default implementation of Logger
type DefaultLogger struct {
	Log *log.Logger
}

// NewDefaultLogger - Create new default logger
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		Log: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (d *DefaultLogger) addLogLevel(level, msg string) string {
	return fmt.Sprintf("[%s] %s", level, msg)
}

func (d *DefaultLogger) addExtraFields(extraFields map[string]interface{}, msg string) string {
	extraString := ""
	for k, v := range extraFields {
		extraString = fmt.Sprintf("%s%s=%v ", extraString, k, v)
	}
	if extraString != "" {
		msg = extraString + msg
	}
	return msg
}

// Debug -
func (d *DefaultLogger) Debug(msg string, extraFields map[string]interface{}) {
	msg = d.addExtraFields(extraFields, msg)
	msg = d.addLogLevel("DEBUG", msg)
	d.Log.Print(msg)
}

// Info -
func (d *DefaultLogger) Info(msg string, extraFields map[string]interface{}) {
	msg = d.addExtraFields(extraFields, msg)
	msg = d.addLogLevel("INFO", msg)
	d.Log.Print(msg)
}

// Warn -
func (d *DefaultLogger) Warn(msg string, extraFields map[string]interface{}) {
	msg = d.addExtraFields(extraFields, msg)
	msg = d.addLogLevel("WARN", msg)
	d.Log.Print(msg)
}

// Error -
func (d *DefaultLogger) Error(msg string, extraFields map[string]interface{}) {
	msg = d.addExtraFields(extraFields, msg)
	msg = d.addLogLevel("ERROR", msg)
	d.Log.Print(msg)
}
