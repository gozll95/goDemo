package client

import (
	"net/http"

	"github.com/qbusapi/lib/httpclient"
	"github.com/qbusapi/lib/token"
)

type SenderAgent interface {
	SetHeaders(key, value string)
	GetHeaders() map[string]string
	GetToken() (token string, err error)
	DoJSONWithHeaders(method string, url string, headers map[string]string, data interface{}, result interface{}) (*http.Response, error)
}

type SenderClient struct {
	Token   token.Token
	Client  *httpclient.HTTPClient
	Headers map[string]string
}

func NewSenderClient(token token.Token) SenderAgent {
	sender := &SenderClient{
		Token:   token,
		Client:  httpclient.NewHTTPClient(nil),
		Headers: make(map[string]string),
	}

	return sender
}

func (s *SenderClient) SetHeaders(key, value string) {
	s.Headers[key] = value
}

func (s *SenderClient) GetToken() (token string, err error) {
	return s.Token.GetToken()
}

func (s *SenderClient) GetHeaders() map[string]string {
	return s.Headers
}

func (s *SenderClient) DoJSONWithHeaders(method string, url string, headers map[string]string, data interface{}, result interface{}) (*http.Response, error) {
	return s.Client.DoJSONWithHeaders(method, url, headers, data, result)
}
