package main

import (
	"fmt"
	"net/http"
	"os"

	logmw "go-kit-gin/middleware/log"
	"go-kit-gin/router"
	"go-kit-gin/service"

	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
	//go log.RenameLogFile()

	// Only log the warning severity or above.
	//logrus.SetLevel(logrus.WarnLevel)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {

	var svc service.AppService
	svc = &service.AppSvc{}
	svc = logmw.LoggingMiddleware()(svc)
	router.Start(svc)

}

func consulCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "consulCheck")
}
