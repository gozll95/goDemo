package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func main() {
	client := &amp
	https.Client{}
	resp, err := client.Get("https://i.cqut.edu.cn/zfca?yhlx=student&amp;login=0122579031373493728&amp;url=xs_main.aspx")
	CheckError(err)
	defer resp.Body.Close()

	//获得cookie
	cookies := resp.Cookies()
	//获得JsessionId，之后会用到
	JsessionId := strings.Replace(cookies[0].String(), "; Path=/zfca", "", 1)
	JsessionId = strings.Trim(JsessionId, "")
	fmt.Println(JsessionId)
	//解析html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	CheckError(err)
	ltnode := doc.Find(`[name="lt"]`)
	lt, exist := ltnode.Attr("value")
	if exist {
		fmt.Println(lt)
	}
	form := url.Values{
		"useValidateCode": {"0"},
		"isremenberme":    {"0"},
		"lt":              {lt},      //第二步解析出啊的lt
		"username":        {"****"},  //你的账号
		"password":        {"*****"}, //你的密码
		"_eventId":        {"submit"},
	}

	formdata := ioutil.NopCloser(strings.NewReader(form.Encode())) //转换成io.reader
	//post的url其实不是固定的，但不要解析之前获得的html，其它部分都是一样的，只有JsessionId，以为我们之前的get是没有cookies的，所以
	//这里要加JsessionId
	urlpost := "https://i.cqut.edu.cn/zfca/login;" + JsessionId + "?yhlx=student&amp;login=0122579031373493728&amp;url=xs_main.aspx"
	//POST是大写，吃过这个亏
	post, err := https.NewRequest("POST", urlpost, formdata)
	//必须指定Content-Type
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//设置cookie
	post.Header.Set("Cookie", JsessionId)
	fmt.Printf("%+v\n", post)
	resp1, err := client.Do(post)
	CheckError(err)
	defer resp1.Body.Close()
	body, err := ioutil.ReadAll(resp1.Body)
	//打印body看是否成功
	fmt.Println(string(body))
}
