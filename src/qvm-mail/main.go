package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"qvm-mail/client"
	"qvm-mail/mail"

	"github.com/astaxie/beego/logs"
)

var (
	QbusHost = "https://127.0.0.1:9090/v1.0"
	QbusAk   = "OTM3ZWI4MzEtZmIzNi00NTg4LWJlNjQtYzYwNWQ3ZWNkY2ZkOlFWTQ"
	QbusSk   = "97b22f8a83cddc88e570a4739d09720e9eaf96c60bcbbc15c4365778e1df1c78"
)

func init() {
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
}

// 需要判断user是否是isActive
func main() {
	Mail()
	//SMS()
	//T()
	//Qbus()
}

func SMS() {
	templateData := mail.TemplateData{
		SMS: &mail.SMSInfo{},
	}
	templateData.SMS.Phone = "18251826682"
	databyte, err := json.Marshal(templateData)
	if err != nil {
		logs.Error(err)
		return
	}
	fmt.Println(string(databyte))
	sendParams := &mail.MailSendParam{
		Host:       QbusHost,
		AccessKey:  QbusAk,
		SecretKey:  QbusSk,
		Data:       string(databyte),
		TemplateId: "5a655880b79f52081946d489",
	}
	err = mail.SendMail(sendParams)
	if err != nil {
		logs.Error(err)
		return
	}
}

func Mail() {
	templateData := mail.TemplateData{
		User:  &mail.UserInfo{},
		User1: &mail.User1Info{},
	}
	templateData.User.Email = "xxxxxxxxxx@bbbb.com"
	templateData.User1.A = "i am xxxxxxxxxx"
	databyte, err := json.Marshal(templateData)
	if err != nil {
		logs.Error(err)
		return
	}
	fmt.Println(string(databyte))
	sendParams := &mail.MailSendParam{
		Host:       QbusHost,
		AccessKey:  QbusAk,
		SecretKey:  QbusSk,
		Data:       string(databyte),
		TemplateId: "5a7806caa9af7e6bc7233783",
	}
	err = mail.SendMail(sendParams)
	if err != nil {
		logs.Error(err)
		return
	}
}

func T() {
	templateData := mail.User1Info{}
	templateData.A = "i am xxxxxxxxxx test"
	databyte, err := json.Marshal(templateData)
	if err != nil {
		logs.Error(err)
		return
	}
	fmt.Println(string(databyte))
	sendParams := &mail.MailSendParam{
		Host:       QbusHost,
		AccessKey:  QbusAk,
		SecretKey:  QbusSk,
		Data:       string(databyte),
		TemplateId: "xxxx",
	}
	err = mail.SendMail(sendParams)
	if err != nil {
		logs.Error(err)
		return
	}
}

// test qbus api
func Qbus() {

	url := QbusHost + "/template"
	request, err := http.NewRequest("GET", url, nil)
	client := client.New(QbusAk, QbusSk)
	resp, err := client.RoundTrip(request)
	if err != nil {
		logs.Error(err)
		return
	}
	fmt.Println(resp)

	if resp.StatusCode != http.StatusOK {
		logs.Error("wrong")
	}
	return
}
