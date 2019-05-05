package main

import (
	"os"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	// You could set this to any `io.Writer` such as a file
	errFile, err := os.OpenFile("logrus-err.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	infoFile, err := os.OpenFile("logrus-info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	// pathMap := lfshook.PathMap{
	// 	logrus.InfoLevel:  "./info.log",
	// 	logrus.ErrorLevel: "./info.log",
	// }

	// 	Log.Hooks.Add(lfshook.NewHook(
	// 	pathMap,
	// 	&logrus.JSONFormatter{},
	// ))

	Log = logrus.New()

	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  infoFile,
		logrus.ErrorLevel: errFile,
	}

	Log.Hooks.Add(lfshook.NewHook(
		writerMap,
		&logrus.JSONFormatter{},
	))

	return Log
}

func main() {
	Log = NewLogger()
	Log.Println("info")
	Log.Error("error2")
}
