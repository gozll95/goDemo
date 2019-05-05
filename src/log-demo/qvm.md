package utils

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Logger struct {
	LogEntry *log.Entry
}

func NewLogger(entry *log.Entry) *Logger {
	return &Logger{
		LogEntry: entry,
	}
}

func NewEmptyLogger() *Logger {
	logger := log.New()

	return &Logger{
		LogEntry: log.NewEntry(logger),
	}
}

var StdLog = NewEmptyLogger()

func (logger *Logger) Debug(args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Debug(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Info(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Error(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Fatal(args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Panic(args...)
}

// Entry Printf family functions
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Debugf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Infof(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Warnf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Errorf(format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Fatalf(format, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.LogEntry.WithFields(log.Fields{"timedate": time.Now()}).Panicf(format, args...)
}
