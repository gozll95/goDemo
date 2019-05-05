package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"qvm-mail/mail"

	mns_params "github.com/zhu/qvm/server/lib/aliyun/mns/params"
	"github.com/zhu/qvm/server/utils"
	"github.com/zhu/qvm/server/utils/multi_tag"
)

const (
	Tag = "ason"
	// QbusHost = "https://qbus.bbbb.com/v1.0"
	// QbusAk   = "OTM3ZWI4MzEtZmIzNi00NTg4LWJlNjQtYzYwNWQ3ZWNkY2ZkOlFWTQ"
	// QbusSk   = "97b22f8a83cddc88e570a4739d09720e9eaf96c60bcbbc15c4365778e1df1c78"

	QbusHost = "http://127.0.0.1:9090/v1.0"
	QbusAk   = "NjkwZGUwZjAtNGQxMy00YTE3LWIyMDktMjhkZDg3MWYzOWI5OlRFU1Q"
	QbusSk   = "b6293367788e2aa3ac389f0f090be6af65de2668781c362cb3e9b3f6787e01cd"
)

type EventType struct {
	UserID   string `ason:"userID"`
	EventID  string `ason:"eventID"`
	SendData string `ason:"-"`
	//Level    mns_params.SendLevel `ason:"-"`
}

type UserInfo struct {
	Email       string `ason:"email"`
	Quid        uint32 `ason:"Quid"`
	UserName    string `ason:"userName"`
	PhoneNumber string `ason:"phone_number"`
}

type MnsMessage struct {
	User *UserInfo   `ason:"user"`
	Msg  interface{} `ason:"msg"`
}

type MailSendParam struct {
	Host       string           `json:"host"`
	AccessKey  string           `json:"access_key"`
	SecretKey  string           `json:"secret_key"`
	TemplateId string           `json:"template_id"`
	Template   *MessageTemplate `json:"template"`
	Data       string           `json:"data"`
}

type MessageTemplate struct {
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	RenderType  string `json:"render_type"`
}

func parseBody(msg string) (event EventType, err error) {
	err = multi_tag.Unmarshal([]byte(msg), &event, Tag)
	if err != nil {
		return event, err
	}

	msgBody, ok := mns_params.ToStruct[event.EventID]
	if !ok {
		return event, errors.New(fmt.Sprintf("params.ToLevel[%v]:no such type", event.EventID))
	}
	err = multi_tag.Unmarshal([]byte(msg), &msgBody, Tag)
	if err != nil {
		utils.StdLog.Errorf("multi_tag.Unmarshal([]byte(%v), %+v, %v):%v", msg, msgBody, Tag, err)
		return event, err
	}

	//toDo: from base account api to make email and phone and username
	user := &UserInfo{
		Email:       "xxxxx",
		Quid:        1381084496, //这里是bbbb user id
		UserName:    "朱立蕾",
		PhoneNumber: "18251826682",
	}
	u := MnsMessage{
		User: user,
		Msg:  msgBody,
	}

	databyte, err := multi_tag.Marshal(u, Tag)
	if err != nil {
		return event, err
	}

	var out bytes.Buffer
	err = json.Indent(&out, databyte, "", "\t")

	if err != nil {
		log.Fatalln(err)
	}

	out.WriteTo(os.Stdout)

	//{1103909446200972 ecs_upgrade {"email":"xx","phone":"xx","userName":"xx","msgInfo":{"data":{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_upgrade","messageSource":"console","source":"console","timeStamp":1516800705630,"timestamp":1516800741,"uniqueID":"3aa1fcea-7dce-4ef0-a28e-c6425619e36c","userID":"1103909446200972"}} 1}

	// msgLevel, ok := mns_params.ToLevel[event.EventID]
	// if !ok {
	// 	return event, errors.New(fmt.Sprintf("params.ToLevel[%v]:no such type", event.EventID))
	// }

	event.SendData = string(databyte)
	//event.Level = msgLevel
	//beego.Error(event)
	return event, nil
}

func send() {

}

//根据event.EventId去数据库选择不同的模板ID
func (event *EventType) SendMail(gId int) (err error) {
	var res mail.Response

	sendParams := &mail.MailSendParam{
		Host:       QbusHost,
		AccessKey:  QbusAk,
		SecretKey:  QbusSk,
		Data:       event.SendData,
		TemplateId: "5a7806caa9af7e6bc7233783",
	}

	qbusClient := mail.NewQbusClient(QbusAk, QbusSk)

	_, err = qbusClient.SendMail(sendParams, &res)
	if err != nil {
		fmt.Printf("send_message.SendMail(%v,%v):%v\n", sendParams, gId, err)
		return
	}
	fmt.Println(res)

	return
}

func (event *EventType) SendMailBybbbb(gId int) (err error) {
	var res mail.Response

	provider := "bbbb"
	sendParams := &mail.MailSendParam{
		Host:       QbusHost,
		AccessKey:  QbusAk,
		SecretKey:  QbusSk,
		Data:       event.SendData,
		TemplateId: "5b3270156b3e7cccae21011e",
		Provider:   &provider,
	}

	qbusClient := mail.NewQbusClient(QbusAk, QbusSk)

	_, err = qbusClient.SendMail(sendParams, &res)
	if err != nil {
		fmt.Printf("send_message.SendMail(%v,%v):%v\n", sendParams, gId, err)
		return
	}
	fmt.Println(res)

	return
}

func (event *EventType) SendSMS(gId int) (err error) {
	var res mail.Response

	sendParams := &mail.MailSendParam{
		Host:       QbusHost,
		AccessKey:  QbusAk,
		SecretKey:  QbusSk,
		Data:       event.SendData,
		TemplateId: "5a782bc9a9af7ea1149a3720",
	}

	qbusClient := mail.NewQbusClient(QbusAk, QbusSk)
	_, err = qbusClient.SendMail(sendParams, &res)
	if err != nil {
		utils.StdLog.Errorf("send_message.SMS(%v,%v):%v", sendParams, gId, err)
		return
	}

	fmt.Println(res)

	messageShowParam := &mail.MessageShowParam{
		Host:      QbusHost,
		AccessKey: QbusAk,
		SecretKey: QbusSk,
		MessageId: fmt.Sprintf("%v", res.Data),
	}
	_, err = qbusClient.ShowMesage(messageShowParam, &res)
	if err != nil {
		fmt.Println(err)
		//	utils.StdLog.Error(err)
		return
	}

	var message mail.MessageModel

	b, err := json.Marshal(res.Data)
	err = json.Unmarshal(b, &message)

	fmt.Println(message.Status)

	return
}

func (event *EventType) SendSMSByMorse(gId int) (err error) {
	var res mail.Response

	provider := "morse"
	sendParams := &mail.MailSendParam{
		Host:       QbusHost,
		AccessKey:  QbusAk,
		SecretKey:  QbusSk,
		Data:       event.SendData,
		TemplateId: "5b327a966b3e7cccae21038c",
		Provider:   &provider,
	}

	qbusClient := mail.NewQbusClient(QbusAk, QbusSk)
	_, err = qbusClient.SendMail(sendParams, &res)
	if err != nil {
		utils.StdLog.Errorf("send_message.SMS(%v,%v):%v", sendParams, gId, err)
		return
	}

	fmt.Println(res)

	messageShowParam := &mail.MessageShowParam{
		Host:      QbusHost,
		AccessKey: QbusAk,
		SecretKey: QbusSk,
		MessageId: fmt.Sprintf("%v", res.Data),
	}
	_, err = qbusClient.ShowMesage(messageShowParam, &res)
	if err != nil {
		fmt.Println(err)
		//	utils.StdLog.Error(err)
		return
	}

	var message mail.MessageModel

	b, err := json.Marshal(res.Data)
	err = json.Unmarshal(b, &message)

	fmt.Println(message.Status)

	return
}

func ShowStatus() {
	var res mail.Response

	messageShowParam := &mail.MessageShowParam{
		Host:      QbusHost,
		AccessKey: QbusAk,
		SecretKey: QbusSk,
		MessageId: "5a7ac774b79f522b0bca30dd",
	}

	qbusClient := mail.NewQbusClient(QbusAk, QbusSk)
	_, err := qbusClient.ShowMesage(messageShowParam, &res)
	if err != nil {
		fmt.Println(err)
		//	utils.StdLog.Error(err)
		return
	}

	var message mail.MessageModel

	b, err := json.Marshal(res.Data)
	err = json.Unmarshal(b, &message)

	fmt.Println(message.Status)
}

func main() {
	//ecs upgrade msg := `{"data":{"instanceId":"i-xxxxxxxxxxxxxxxxxxxx","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_upgrade","messageSource":"console","source":"console","timeStamp":1516800705630,"timestamp":1516800741,"uniqueID":"3aa1fcea-7dce-4ef0-a28e-c6425619e36c","userID":"1103909446200972"}`
	//msg := `{"data":{"diskCount":"7","endTime":"2016-12-22 10:10","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","keepDays":"8","lastDay":"2016-12-22 10:10","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_expired","messageSource":"console","source":"console","timeStamp":1516876173680,"timestamp":1516876176,"uniqueID":"5aeaa9ad-e43e-44d6-afd8-246e4cd65927","userID":"1103909446200972"}`
	event, err := parseBody(ecs_upgrade_msg)
	if err != nil {
		panic(err)
	}
	// err = event.SendMail(1)
	// if err != nil {
	// 	panic(err)
	// }
	// err = event.SendSMS(1)
	// if err != nil {
	// 	panic(err)
	// }

	err = event.SendMailBybbbb(1)
	if err != nil {
		panic(err)
	}

	//ShowStatus()
}
