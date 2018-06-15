package log

var stdLog Logger

// Logger - Interface used by taskor to log
type Logger interface {
	Debug(msg string, extraFields map[string]interface{})
	Info(msg string, extraFields map[string]interface{})
	Warn(msg string, extraFields map[string]interface{})
	Error(msg string, extraFields map[string]interface{})
}

func init() {
	stdLog = NewDefaultLogger()
}

// SetLogger - change current logger
func SetLogger(newLogger Logger) {
	stdLog = newLogger
}

// Debug log with level debug
func Debug(msg string) {
	stdLog.Debug(msg, nil)
}

// DebugWithFields log with extraFields with level debug
func DebugWithFields(msg string, i interface{}) {
	stdLog.Debug(msg, GetFields(i))
}

// Info log with level Info
func Info(msg string) {
	stdLog.Info(msg, nil)
}

// InfoWithFields log with extraFields with level Info
func InfoWithFields(msg string, i interface{}) {
	stdLog.Info(msg, GetFields(i))
}

// Warn log with level Warn
func Warn(msg string) {
	stdLog.Warn(msg, nil)
}

// WarnWithFields log with extraFields with level Warn
func WarnWithFields(msg string, i interface{}) {
	stdLog.Warn(msg, GetFields(i))
}

// Error log with level Error
func Error(msg string) {
	stdLog.Error(msg, nil)
}

// ErrorWithFields log with extraFields with level Error
func ErrorWithFields(msg string, i interface{}) {
	stdLog.Error(msg, GetFields(i))
}
