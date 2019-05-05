package utils

import (
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	LogEntry *logrus.Entry
}

func NewLogger(entry *logrus.Entry) *Logger {
	return &Logger{
		LogEntry: entry,
	}
}

func NewEmptyLogger() *Logger {
	logger := logrus.New()

	return &Logger{
		LogEntry: logrus.NewEntry(logger),
	}
}

var (
	StdLog  = NewLogger(DecorateRuntimeContext(logrus.NewEntry(logrus.New())))
	StdLog2 = NewEmptyLogger()
)

func (logger *Logger) Debug(args ...interface{}) {
	logger.LogEntry.WithFields(logrus.Fields{"timedate": time.Now()}).Debug(args...)
}

// func (logger *Logger) Info(args ...interface{}) {
// 	logger.LogEntry.WithFields(logrus.Fields{"timedate": time.Now()}).Info(args...)
// }

func (logger *Logger) Info(args ...interface{}) {
	DecorateRuntimeContext(logger.LogEntry.WithFields(logrus.Fields{"timedate": time.Now()})).Info(args...)
}

func DecorateRuntimeContext(logger *logrus.Entry) *logrus.Entry {

	// if pc, file, line, ok := runtime.Caller(2); ok {
	// 	fName := runtime.FuncForPC(pc).Name()
	// 	return logger.WithField("file", file).WithField("line", line).WithField("func", fName)
	// } else {
	// 	return logger
	// }

	//	_, file, line, ok := runtime.Caller(l.callDepth)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	fileSlice := strings.Split(file, "log-demo")

	a := filepath.Join("log-demo", fileSlice[1]) + ":" + itoa(line, -1)
	return logger.WithField("file", a)
}

func itoa(i int, wid int) string {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		return "0"
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	return string(b[bp:])
}
