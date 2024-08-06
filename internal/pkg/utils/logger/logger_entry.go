package logger

import "github.com/sirupsen/logrus"

type LoggerEntry struct {
	entry *logrus.Entry
}

func (logger *LoggerEntry) Errorf(format string, args ...any) {
	logger.entry.Errorf(format, args...)
}

func (logger *LoggerEntry) Error(args ...any) {
	logger.entry.Error(args...)
}

func (logger *LoggerEntry) Fatalf(format string, args ...any) {
	logger.entry.Fatalf(format, args...)
}

func (logger *LoggerEntry) Fatal(args ...any) {
	logger.entry.Fatal(args...)
}

func (logger *LoggerEntry) Warn(args ...any) {
	logger.entry.Warn(args...)
}

func (logger *LoggerEntry) Warnf(format string, args ...any) {
	logger.entry.Warnf(format, args...)
}

func (logger *LoggerEntry) Debugf(format string, args ...any) {
	logger.entry.Debugf(format, args...)
}

func (logger *LoggerEntry) Debug(args ...any) {
	logger.entry.Debug(args...)
}

func (logger *LoggerEntry) Infof(format string, args ...any) {
	logger.entry.Infof(format, args...)
}

func (logger *LoggerEntry) Info(args ...any) {
	logger.entry.Info(args...)
}

func (logger *LoggerEntry) WithField(key string, value any) (entry Logger) {
	entry = &LoggerEntry{logger.entry.WithField(key, value)}
	return
}

func (logger *LoggerEntry) WithFields(args map[string]any) (entry Logger) {
	entry = &LoggerEntry{logger.entry.WithFields(args)}
	return
}
