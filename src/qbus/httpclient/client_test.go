package httpclient

import (
	"net/http"
	"testing"

	"github.com/dolab/httpmitm"
	"github.com/golib/assert"
)

func Test_HTTPClient(t *testing.T) {
	mitm := httpmitm.NewMitmTransport().StubDefaultTransport(t)
	defer mitm.UnstubDefaultTransport()

	var (
		url       = "http://www.example.com"
		assertion = assert.New(t)
	)

	mitm.MockRequest("GET", url).WithJsonResponse(http.StatusOK, nil, `{"ok":true}`)

	client := NewHTTPClient(&http.Client{
		Transport: mitm,
	})

	resp, err := client.Do("GET", "mitm"+url[4:], nil)
	assertion.Nil(err)
	assertion.Equal(http.StatusOK, resp.StatusCode)
}

func Test_HTTPClient_DoJSON(t *testing.T) {
	mitm := httpmitm.NewMitmTransport().StubDefaultTransport(t)
	defer mitm.UnstubDefaultTransport()

	var (
		url       = "http://www.example.com"
		assertion = assert.New(t)
	)

	mitm.MockRequest("POST", url).WithJsonResponse(http.StatusOK, nil, `{"ok":true}`)

	client := NewHTTPClient(&http.Client{
		Transport: mitm,
	})

	resp, err := client.DoJSON("POST", "mitm"+url[4:], nil, nil)
	assertion.Nil(err)
	assertion.Equal(http.StatusOK, resp.StatusCode)
}

func Test_HTTPClient_DoXML(t *testing.T) {
	mitm := httpmitm.NewMitmTransport().StubDefaultTransport(t)
	defer mitm.UnstubDefaultTransport()

	var (
		url       = "http://www.example.com"
		assertion = assert.New(t)
	)

	mitm.MockRequest("POST", url).WithXmlResponse(http.StatusOK, nil, `<xml><ok>true</ok></xml>`)

	client := NewHTTPClient(&http.Client{
		Transport: mitm,
	})

	resp, err := client.DoXML("POST", "mitm"+url[4:], nil, nil)
	assertion.Nil(err)
	assertion.Equal(http.StatusOK, resp.StatusCode)
}
