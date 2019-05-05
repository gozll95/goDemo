package httpclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// UserClient is user independent http request, and each of the request has a context
// Note: UserClient is not safe for concurrent use by multiple goroutines
type UserClient struct {
	client *http.Client

	response *http.Response
	rawBody  []byte
	rawError error
	isReaded bool

	host    string
	header  http.Header
	isHttps bool
}

func NewUserClient(host string, isHttps bool) *UserClient {
	// adjust host
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		purl, err := url.Parse(host)
		if err == nil {
			host = purl.Host
			isHttps = (purl.Scheme == "https")
		}
	}

	jar, _ := cookiejar.New(nil)

	return &UserClient{
		client: &http.Client{
			Jar: jar,
		},
		host:    host,
		isHttps: isHttps,
	}

}

// Host returns the host and port of the server, e.g. "127.0.0.1:9090"
func (uc *UserClient) Host() string {
	switch {
	case uc.host == "":
		return "127.0.0.1"

	case uc.host[0] == ':':
		return "127.0.0.1" + uc.host

	}

	return uc.host
}

// Abs returns the abs http/https URL of the resource, e.g. "http://127.0.0.1:9090/status".
// The scheme is set to https if client created with isHttps == true.
func (uc *UserClient) Abs(path string) string {
	scheme := "http://"
	if uc.isHttps {
		scheme = "https://"
	}

	return scheme + uc.Host() + path
}

// Cookies returns cookies related to the host
func (uc *UserClient) Cookies() []*http.Cookie {
	purl, _ := url.Parse(uc.Abs("/"))

	return uc.client.Jar.Cookies(purl)
}

// SetCookie sets cookies with the host
func (uc *UserClient) SetCookies(cookies []*http.Cookie) {
	purl, _ := url.Parse(uc.Abs("/"))

	uc.client.Jar.SetCookies(purl, cookies)
}

// Headers returns default headers
func (uc *UserClient) Headers() http.Header {
	return uc.header
}

// SetHeader sets default header of the request
func (uc *UserClient) SetHeader(key, value string) {
	uc.header.Set(key, value)
}

// NewRequest issues any request with injected default header
// If successful, the caller may examine the client.Response() properties.
// NOTE: You have to manage session / cookie data manually.
func (uc *UserClient) NewRequest(request *http.Request) (err error) {
	// inject default header
	for key, value := range uc.header {
		if _, ok := request.Header[key]; !ok {
			request.Header[key] = value
		}
	}

	uc.response, err = uc.client.Do(request)
	uc.rawBody = []byte{}
	uc.rawError = nil
	uc.isReaded = false

	return
}

// NewSessionRequest issues any request with injected session / cookie
// If successful, the caller may examine the client.Response() properties.
// NOTE: Session data will be added to the request cookies for you.
func (uc *UserClient) NewSessionRequest(request *http.Request) error {
	// inject releated cookies
	for _, cookie := range uc.client.Jar.Cookies(request.URL) {
		request.AddCookie(cookie)
	}

	return uc.NewRequest(request)
}

// NewMultipartRequest issues a multipart request for the method & fields given
// If successful, the caller may examine the client.Response() properties.
// NOTE: Session data will be added to the request cookies for you.
func (uc *UserClient) NewMultipartRequest(method, path, filename string, file interface{}, fields ...map[string]string) error {
	var buf bytes.Buffer

	mw := multipart.NewWriter(&buf)

	fw, ferr := mw.CreateFormFile("filename", filename)
	if ferr != nil {
		return ferr
	}

	// apply file
	var (
		reader io.Reader
		err    error
	)
	switch file.(type) {
	case io.Reader:
		reader, _ = file.(io.Reader)

	case *os.File:
		reader, _ = file.(*os.File)

	case string:
		filepath, _ := file.(string)

		reader, err = os.Open(filepath)
		if err != nil {
			return err
		}

	}

	if _, err := io.Copy(fw, reader); err != nil {
		return err
	}

	// apply fields
	if len(fields) > 0 {
		for key, value := range fields[0] {
			mw.WriteField(key, value)
		}
	}

	// adds the terminating boundary
	mw.Close()

	request, err := http.NewRequest(method, uc.Abs(path), &buf)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", mw.FormDataContentType())

	return uc.NewSessionRequest(request)
}

// Response returns *http.Response of the previous request
func (uc *UserClient) Response() *http.Response {
	return uc.response
}

// RawBody returns response body of the previous request
// It read response.Body if it has not been read.
func (uc *UserClient) RawBody() ([]byte, error) {
	if uc.isReaded {
		return uc.rawBody, uc.rawError
	}
	uc.isReaded = true

	defer uc.response.Body.Close()

	if uc.response.StatusCode < 200 || uc.response.StatusCode > 300 {
		uc.rawError = errors.New("Unexpected response status code:" + strconv.Itoa(uc.response.StatusCode) + ":" + uc.response.Status)
	} else {
		uc.rawBody, uc.rawError = ioutil.ReadAll(uc.response.Body)
	}

	return uc.rawBody, uc.rawError
}

// StreamBody copy response.Body of the previous request to target io.Writer
// It useful for reading response.Body with big bytes.
// NOTE: It panics with ErrResponseRead if the response.Body has been read.
func (uc *UserClient) StreamBody(w io.Writer) error {
	if uc.isReaded {
		panic(ErrResponseRead.Error())
	}
	uc.isReaded = true

	defer uc.response.Body.Close()

	_, uc.rawError = io.Copy(w, uc.response.Body)

	return uc.rawError
}

// UnmarshalJSON read response.Body and decode data with json.Unmarshal
func (uc *UserClient) UnmarshalJSON(v interface{}) error {
	data, err := uc.RawBody()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &v)
}

// UnmarshalXML read response.Body and decode data with xml.Unmarshal
func (uc *UserClient) UnmarshalXML(v interface{}) error {
	data, err := uc.RawBody()
	if err != nil {
		return err
	}

	return xml.Unmarshal(data, &v)
}

// Discard close response.Body without reading
func (uc *UserClient) Discard() error {
	if uc.isReaded {
		return uc.rawError
	}
	uc.isReaded = true

	uc.rawError = uc.response.Body.Close()

	return uc.rawError
}

// Get issues a GET request to the given path and stores the result in client.Response().
// NOTE: Session data will be added to the request cookies for you.
func (uc *UserClient) Get(path string, params ...url.Values) error {
	lru := uc.Abs(path)
	lru = Client.UrlEncode(lru, params...)

	request, err := http.NewRequest("GET", lru, nil)
	if err != nil {
		return err
	}

	return uc.NewSessionRequest(request)
}

// Head issues a HEAD request to the given path and stores the result in client.Response().
// NOTE: Session data will be added to the request cookies for you.
func (uc *UserClient) Head(path string, params ...url.Values) error {
	lru := uc.Abs(path)
	lru = Client.UrlEncode(lru, params...)

	request, err := http.NewRequest("HEAD", lru, nil)
	if err != nil {
		return err
	}

	return uc.NewSessionRequest(request)
}

// Options issues an OPTIONS request to the given path and stores the result in client.Response().
// NOTE: Session data will be added to the request cookies for you.
func (uc *UserClient) Options(path string, params ...url.Values) error {
	lru := uc.Abs(path)
	lru = Client.UrlEncode(lru, params...)

	request, err := http.NewRequest("OPTIONS", lru, nil)
	if err != nil {
		return err
	}

	return uc.NewSessionRequest(request)
}

// Put issues a PUT request to the given path, sending request with specified Content-Type header, and
// stores the result in client.Response().
func (uc *UserClient) Put(path, contentType string, data ...interface{}) error {
	return uc.Invoke("PUT", path, contentType, data...)
}

// PutForm issues a PUT request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in client.Response().
func (uc *UserClient) PutForm(path string, data interface{}) error {
	return uc.Put(path, "application/x-www-form-urlencoded", data)
}

// PutJSON issues a PUT request to the given path with Content-Type: application/json header, and
// stores the result in client.Response().
// It will encode data by json.Marshal before making request.
func (uc *UserClient) PutJSON(path string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return uc.Put(path, "application/json", b)
}

// PutXML issues a PUT request to the given path with Content-Type: text/xml header, and
// stores the result in client.Response().
// It will encode data by xml.Marshal before making request.
func (uc *UserClient) PutXML(path string, data interface{}) error {
	b, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	return uc.Put(path, "text/xml", b)
}

// Post issues a POST request to the given path, sending request with specified Content-Type header, and
// stores the result in client.Response().
func (uc *UserClient) Post(path, contentType string, data ...interface{}) error {
	return uc.Invoke("POST", path, contentType, data...)
}

// PostForm issues a POST request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in client.Response().
func (uc *UserClient) PostForm(path string, data interface{}) error {
	return uc.Post(path, "application/x-www-form-urlencoded", data)
}

// PostJSON issues a POST request to the given path with Content-Type: application/json header, and
// stores the result in client.Response().
// It will encode data by json.Marshal before making request.
func (uc *UserClient) PostJSON(path string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return uc.Post(path, "application/json", b)
}

// PostXML issues a POST request to the given path with Content-Type: text/xml header, and
// stores the result in client.Response().
// It will encode data by xml.Marshal before making request.
func (uc *UserClient) PostXML(path string, data interface{}) error {
	b, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	return uc.Post(path, "text/xml", b)
}

// Patch issues a PATCH request to the given path, sending request with specified Content-Type header, and
// stores the result in client.Response().
func (uc *UserClient) Patch(path, contentType string, data ...interface{}) error {
	return uc.Invoke("PATCH", path, contentType, data...)
}

// PatchForm issues a PATCH request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in client.Response().
func (uc *UserClient) PatchForm(path string, data interface{}) error {
	return uc.Patch(path, "application/x-www-form-urlencoded", data)
}

// PatchJSON issues a PATCH request to the given path with with Content-Type: application/json header, and
// stores the result in client.Response().
// It will encode data by json.Marshal before making request.
func (uc *UserClient) PatchJSON(path string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return uc.Patch(path, "application/json", b)
}

// PatchXML issues a PATCH request to the given path with Content-Type: text/xml header, and
// stores the result in client.Response().
// It will encode data by xml.Marshal before making request.
func (uc *UserClient) PatchXML(path string, data interface{}) error {
	b, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	return uc.Patch(path, "text/xml", b)
}

// Delete issues a DELETE request to the given path, sending request with specified Content-Type header, and
// stores the result in client.Response().
func (uc *UserClient) Delete(path, contentType string, data ...interface{}) error {
	return uc.Invoke("DELETE", path, contentType, data...)
}

// DeleteForm issues a DELETE request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in client.Response().
func (uc *UserClient) DeleteForm(path string, data interface{}) error {
	return uc.Delete(path, "application/x-www-form-urlencoded", data)
}

// DeleteJSON issues a DELETE request to the given path with Content-Type: application/json header, and
// stores the result in client.Response().
// It will encode data by json.Marshal before making request.
func (uc *UserClient) DeleteJSON(path string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return uc.Delete(path, "application/json", b)
}

// DeleteXML issues a DELETE request to the given path with Content-Type: text/xml header, and
// stores the result in client.Response().
// It will encode data by xml.Marshal before making request.
func (uc *UserClient) DeleteXML(path string, data interface{}) error {
	b, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	return uc.Delete(path, "text/xml", b)
}

// Invoke issues a HTTP request to the given path with specified method and content type header, and
// stores the result in client.Response().
func (uc *UserClient) Invoke(method, path, contentType string, data ...interface{}) error {
	var (
		request *http.Request
		err     error
	)

	if len(data) == 0 {
		request, err = http.NewRequest(method, uc.Abs(path), nil)
	} else {
		var reader io.Reader

		body := data[0]
		switch body.(type) {
		case io.Reader:
			reader, _ = body.(io.Reader)

		case string:
			s, _ := body.(string)

			reader = bytes.NewBufferString(s)

		case []byte:
			buf, _ := body.([]byte)

			reader = bytes.NewBuffer(buf)

		case url.Values:
			params, _ := body.(url.Values)

			reader = bytes.NewBufferString(params.Encode())

		default:
			reader = strings.NewReader(fmt.Sprintf("%v", body))

		}

		request, err = http.NewRequest(method, uc.Abs(path), reader)
	}
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", contentType)

	return uc.NewSessionRequest(request)
}
