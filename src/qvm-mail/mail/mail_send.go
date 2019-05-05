package mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"qvm-mail/client"

	"github.com/astaxie/beego/logs"
)

func init() {
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
}

type QbusClient struct {
	Client *client.Client
}

func NewQbusClient(ak, sk string) *QbusClient {
	return &QbusClient{
		Client: client.New(ak, sk),
	}
}

func (qbus *QbusClient) SendMail(param *MailSendParam, result interface{}) (*http.Response, error) {
	paramByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(paramByte)
	url := param.Host + "/message"
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	return qbus.doReq(request, result)
}

func (qbus *QbusClient) ShowMesage(param *MessageShowParam, result interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/admin/message/%s", param.Host, param.MessageId)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return qbus.doReq(request, result)
}

func (qbus *QbusClient) doReq(request *http.Request, result interface{}) (*http.Response, error) {
	resp, err := qbus.Client.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("doReq(request %+v, result %+v):code: %v", request, result, resp.StatusCode))
	}

	defer resp.Body.Close()

	return resp, json.NewDecoder(resp.Body).Decode(&result)

}

type Response struct {
	statusCode int         `json:"-"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Resource   string      `json:"resource,omitempty"`
	RequestId  string      `json:"request_id,omitempty"`
	Data       interface{} `json:"data"`
}

type MessageModel struct {
	Id      string `bson:"_id" json:"id"`
	AppId   string `bson:"app_id" json:"app_id"`
	Subject string `bson:"subject" json:"subject"`
	From    string `bson:"from" json:"from"`
	To      string `bson:"to" json:"to"`
	UserId  uint32 `bson:"user_id" json:"user_id"`
	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Status  string `bson:"status" json:"status"`
}
