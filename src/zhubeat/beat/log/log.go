package log

import (
	"io/ioutil"
	"time"
	"zhubeat/beat/log/hook"

	log "github.com/sirupsen/logrus"

	"github.com/zhu/qvm/server/utils"
)

var Log *utils.Logger

func init() {

	start := time.Now()
	path := "/hellp"
	sourceIP := "127.0.0.1"
	method := "GET"
	requestID := "xxxxxxxx"

	Log = utils.NewEmptyLogger()

	Log.LogEntry.WithFields(log.Fields{
		"request_id": requestID,
		"time_start": start,
		"method":     method,
		"source_ip":  sourceIP,
		"path":       path,
	})

	utils.SetFormat("json", Log)
	Log.SetOutput(ioutil.Discard)
	Log.AddHook(hook.NewZhuHook())
}
