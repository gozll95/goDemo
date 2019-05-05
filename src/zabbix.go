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
	spider_base_url     string = "http://192.168.20.25/chart2.php?graphid=547&period=60&stime=20180921003035&updateProfile=1&profileIdx=web.screens&profileIdx2=547"
	login_url           string = "http://192.168.20.25/zabbix"
	post_login_info_url string = "http://192.168.20.25/index.php"
	username            string = "Admin"
	password            string = "zabbix"
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
	f, err := os.Create("download.jpg")
	if err != nil {
		panic(err)
	}

	res, _ := test(login_url, "GET", nil, nil)

	var temp_cookies = res.Cookies()
	fmt.Println(temp_cookies[0].Name, temp_cookies[0].Value)

	defer res.Body.Close()

	fmt.Println(gCurCookieJar)

	// postDict := map[string]string{}
	// postDict["sid"] = "9e47f2f15f4b6735"
	// postDict["form_refresh"] = "1"
	// postDict["name"] = username
	// postDict["password"] = password
	// postDict["from"] = "/"
	// postDict["autologin"] = "1"
	// postDict["enter"] = "Sign+in"

	// //验证
	// requestHeader := map[string]string{}
	// requestHeader["Content-Type"] = "application/x-www-form-urlencoded"
	// requestHeader["Host"] = "192.168.20.25"
	// requestHeader["Origin"] = "http://192.168.20.25"
	// requestHeader["Referer"] = "http://192.168.20.25/zabbix"
	// requestHeader["Upgrade-Insecure-Requests"] = "1"
	// requestHeader["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"

	// httpResp, err := test(post_login_info_url, "POST", postDict, requestHeader)

	urlurl := "http://192.168.20.25/index.php?sid=9e47f2f15f4b6735&form_refresh=1&name=Admin&password=zabbix&autologin=1&enter=Sign+in"

	httpResp, err := test(urlurl, "GET", nil, nil)

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
