package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type loggerLogrus struct {
	log *logrus.Logger
}

func (l *loggerLogrus) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *loggerLogrus) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *loggerLogrus) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *loggerLogrus) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *loggerLogrus) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *loggerLogrus) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *loggerLogrus) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *loggerLogrus) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *loggerLogrus) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *loggerLogrus) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *loggerLogrus) WithField(key string, value interface{}) Logger {
	entry := &LoggerEntry{l.log.WithField(key, value)}
	return entry
}

func (l *loggerLogrus) WithFields(args map[string]interface{}) Logger {
	entry := &LoggerEntry{l.log.WithFields(args)}
	return entry
}

func NewLogger() Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		ForceColors:     true,
	})
	log.SetLevel(logrus.InfoLevel)
	
	log.SetOutput(os.Stdout)

	return &loggerLogrus{
		log: log,
	}
}
