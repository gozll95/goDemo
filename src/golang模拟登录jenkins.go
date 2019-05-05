package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var cookies []*http.Cookie

var gCurCookieJar *cookiejar.Jar

const (
)

func initAll() {
	gCurCookieJar, _ = cookiejar.New(nil)
	return
}
func test(strUrl string, method string, postDict map[string]string, requestHeader map[string]string) (*http.Response, error) {
	httpClient := &http.Client{
		Jar: gCurCookieJar,
	}

	fmt.Println(gCurCookieJar)

	var httpReq *http.Request

	switch method {
	case "GET":
		httpReq, _ = http.NewRequest("GET", strUrl, nil)
	case "POST":
		postValues := url.Values{}
		for postKey, PostValue := range postDict {
			postValues.Set(postKey, PostValue)
		}
		postBytesReader := bytes.NewReader([]byte(postValues.Encode()))
		httpReq, _ = http.NewRequest("POST", strUrl, postBytesReader)

		for headerKey, headerValue := range requestHeader {
			httpReq.Header.Add(headerKey, headerValue)
		}
	}

	httpResp, err := httpClient.Do(httpReq)

	return httpResp, err
}

func login() {
	//创建文件
	f, err := os.Create("download.tar.gz")
	if err != nil {
		panic(err)
	}

	res, _ := test(login_url, "GET", nil, nil)

	var temp_cookies = res.Cookies()
	fmt.Println(temp_cookies[0].Name, temp_cookies[0].Value)

	defer res.Body.Close()

	fmt.Println(gCurCookieJar)

	postDict := map[string]string{}
	postDict["j_username"] = username
	postDict["j_password"] = password
	postDict["from"] = "/"
	postDict["json"] = "xxxxx"
	postDict["Submit"] = "xxxx"

	//验证
	requestHeader := map[string]string{}
	requestHeader["Content-Type"] = "application/x-www-form-urlencoded"
	requestHeader["Host"] = "jenkins.xxxx.io"
	requestHeader["Origin"] = "https://jenkins.xxxx.io"
	requestHeader["Referer"] = "https://jenkins.xxxx.io/login?from=/"

	httpResp, err := test(post_login_info_url, "POST", postDict, requestHeader)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("response Status:", httpResp.Status)
	fmt.Println("response Headers:", httpResp.Header)

	fmt.Println("now cookie is", httpResp.Cookies())

	res1, _ := test(spider_base_url, "GET", nil, nil)

	fmt.Println("response Status:", res1.Status)
	fmt.Println("response Headers:", res1.Header)

	io.Copy(f, res1.Body)
	defer res1.Body.Close()
}

func main() {

	initAll()

	login()
}
