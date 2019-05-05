package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/deckarep/golang-set"
	"github.com/zhu/qvm/server/application"
	"github.com/zhu/qvm/server/application/conf"
	"github.com/zhu/qvm/server/lib/send_message"
	send_params "github.com/zhu/qvm/server/lib/send_message/params"
	"github.com/zhu/qvm/server/model"
	"github.com/zhu/qvm/server/utils"
)

var (
	runMode string // app run mode, available values are [development|test|production], default to development
	srcPath string // app source path
)

func setup() {
	flag.StringVar(&runMode, "runMode", "development", "app -runMode=[development|test|production]")
	flag.StringVar(&srcPath, "srcPath", "", "app -srcPath=/path/to/source")
	flag.Parse()

	mode := conf.ModeType(runMode)

	// verify run mode
	if !mode.IsValid() {
		flag.PrintDefaults()
		return
	}

	// adjust src path
	if srcPath == "" {
		var err error

		srcPath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		srcPath = path.Clean(srcPath)
	}
	items := strings.Split(srcPath, "tools")
	srcPath = items[0]

	// set stdlog and set level
	utils.StdLog = utils.NewEmptyLogger()

	setConfig(mode, srcPath)
}

func setConfig(mode conf.ModeType, srcPath string) (app *application.Application, err error) {
	app = &application.Application{}

	// set config
	app.Config, err = conf.NewConfig(mode, srcPath)
	if err != nil {
		panic(err)
	}
	conf.SetupConfig(app.Config)

	// set model
	mongo := model.NewModel(app.Config.Mongo, utils.StdLog)
	model.SetupModel(mongo)

	return
}

type Qbus struct {
	Host      string `json:"host"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

func message() {
	var (
		res     send_params.Response
		message []send_params.MessageModel
		marker  string = ""
		limit   int    = 100

		kide = mapset.NewSet()
	)
	client := send_message.NewQbusClient(conf.AppConfig.Qbus.QbusAk, conf.AppConfig.Qbus.QbusSk)
	messageShowParam := &send_params.MessageShowParam{
		Host:      conf.AppConfig.Qbus.QbusHost,
		AccessKey: conf.AppConfig.Qbus.QbusAk,
		SecretKey: conf.AppConfig.Qbus.QbusSk,
		MessageId: "",
	}

	for {
		_, err := client.IndexMessage(messageShowParam, limit, marker, &res)
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(res.Data)
		if err != nil {
			return
		}

		err = json.Unmarshal(b, &message)
		if err != nil {
			return
		}

		utils.StdLog.Infof("look for %d qbus items", len(message))
		if len(message) == 0 {
			break
		}

		for _, v := range message {
			kide.Add(v.Id)
			if strings.Contains(v.Subject, "qvm_expire") {
				fmt.Printf("id:%v\tsubject:\t%v\tcreate_time:%v\n", v.Id, v.Subject, v.UpdatedAt)
			}
		}

		marker = message[len(message)-1].Id

	}
}

func main() {
	setup()

	message()

}





func (qbus *QbusClient) IndexMessage(param *params.MessageShowParam, limit int, marker string, result interface{}) (*http.Response, error) {
	var url string
	url = fmt.Sprintf("%s/admin/message?limit=%d&marker=%s", param.Host, limit, marker)
	if marker == "" {
		url = fmt.Sprintf("%s/admin/message?limit=%d", param.Host, limit)
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return qbus.doReq(request, result)
}