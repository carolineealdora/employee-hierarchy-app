package logger

var Log Logger

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(map[string]interface{}) Logger
}

func SetLogger(log Logger) {
	Log = log
}
