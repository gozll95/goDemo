// 登陆原理

// 因为http协议是无状态、短连接的，即浏览器发一个请求给目标服务器，然后服务器把数据传回来， 然后浏览器会服务器之间的交互就结束了。

// 所有需要一种机制来保证让服务器确定用户的身份（比如用户是否登陆），这个机制就是session和cookie， 其中session是服务端保存你的数据，而cookie则是保存在浏览器里面。

// 那么登陆是一个什么过程？其实就是比如说，先输入用户名和密码，把数据post给服务器，服务器会对填写的数据进行校验， 如果校验通过，服务器会生成一段session信息，返回给你的浏览器，你的浏览器会把服务器返回来的session保存在本地， 到了浏览器上，也就变成cookie了。比如说登陆一个网站成功了， 服务器在对应用户对应的session里面添加了一个is_login=true的字段。 浏览器在下次访问这个域名的网站下的其他页面时候，会把存的cookie发送给服务器。

// 如何模拟登陆带验证码的网站

// 其实在打开带验证码的页面的时候，我们就相当于有一次“登陆”，不过这个登陆不是为我们的cookie添加 一个is_login字段，而是添加一个verify-code或者叫captcha字段，服务器会把验证码存在里面，这样它 才知道它生成的每个验证码，是为了那个用户生成的。

// 所以模拟登陆待验证码的网站需要有以下几个步骤：

// 1、访问需要登陆的页面，并获取cookie
// 2、带着cookie，下载验证码图片
// 3、带着cookie，把用户数据（用户名，密码，验证码）通过post方式发给服务器，并且拿到cookie。
// 经过以上几个步骤后，就可以拿到登陆的cookie了，可以带着这个cookie去访问需要权限的页面。 以下是实现代码，go语言编写。(这段代码会首先模拟登陆极客学院，然后爬取职业路径中docker的视频下载链接， 然后生成一堆wget组成的shell脚本，运行脚本就可以下载视频，不过仅对课程少的视频管用， 极客的视频播放的地址会时时刻刻变更)

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

var cookies []*http.Cookie

const (
	spider_base_url     string = "http://www.jikexueyuan.com/path/docker/"
	login_url           string = "http://passport.jikexueyuan.com/sso/login"
	verify_code_url     string = "http://passport.jikexueyuan.com/sso/verify"
	post_login_info_url string = "http://passport.jikexueyuan.com/submit/login?is_ajax=1"
	username            string = "xxxxxxxxxx"
	password            string = "xxxxxxxxxx1995"
)

func getResultHtml(get_url string) *http.Response {
	c := &http.Client{}
	Jar, _ := cookiejar.New(nil)
	getURL, _ := url.Parse(get_url)
	Jar.SetCookies(getURL, cookies)
	c.Jar = Jar
	res, _ := c.Get(get_url)
	return res
}
func login() {
	//获取登陆界面的cookie
	c := &http.Client{}
	req, _ := http.NewRequest("GET", login_url, nil)
	res, _ := c.Do(req)
	req.URL, _ = url.Parse(verify_code_url)
	var temp_cookies = res.Cookies()
	fmt.Println("temp_cookies is", temp_cookies)
	for _, v := range res.Cookies() {
		req.AddCookie(v)
	}
	// 获取验证码
	var verify_code string
	for {
		res, _ = c.Do(req)
		file, _ := os.Create("verify.gif")
		io.Copy(file, res.Body)
		//最开始测试没问题，现在好像改了，
		//最好一次输入正确，然后10秒以内输入，
		//否则，会报出connection reset by peer
		fmt.Println("请查看verify.gif， 然后输入验证码， 看不清输入0重新获取验证码")
		fmt.Scanf("%s", &verify_code)
		if verify_code != "0" {
			break
		}
		res.Body.Close()
	}
	//post数据
	postValues := url.Values{}
	postValues.Add("expire", "7")
	postValues.Add("referer", "http%3A%2F%2Fwww.jikexueyuan.com%2F")
	postValues.Add("uname", username)
	postValues.Add("password", password)
	postValues.Add("verify", verify_code)
	postURL, _ := url.Parse(post_login_info_url)
	Jar, _ := cookiejar.New(nil)
	Jar.SetCookies(postURL, temp_cookies)
	c.Jar = Jar
	res, _ = c.PostForm(post_login_info_url,
		postValues)
	cookies = res.Cookies()
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Println(string(data))
}

type DownloadBean struct {
	dirname  string
	filename string
	href     string
}

func main() {
	login()
	for _, v := range cookies {
		fmt.Println(v)
	}
	var bash_str string = "#!/bin/bash \n"
	baseMap := map[string]string{}
	doc, _ := goquery.NewDocumentFromResponse(getResultHtml(spider_base_url))
	doc.Find(".lesson-info-h2").Each(func(i int, contentSelection *goquery.Selection) {
		selection := contentSelection.Find("a")
		base_href, _ := selection.Attr("href")
		dir_name := selection.Text()
		bash_str += "mkdir \"" + dir_name + "\"\n"
		baseMap[dir_name] = base_href
		fmt.Println(dir_name, "-->", base_href)
	})
	downloadList := []DownloadBean{}
	for k, v := range baseMap {
		doc, _ := goquery.NewDocumentFromResponse(getResultHtml(v))
		doc.Find(".lessonvideo-list dd h2").Each(func(i int, contentSelection *goquery.Selection) {
			selection := contentSelection.Find("a")
			href, _ := selection.Attr("href")
			filename := selection.Text()
			fmt.Println(k, "-->", filename, "-->", href)
			bean := DownloadBean{dirname: k, href: href, filename: filename}
			downloadList = append(downloadList, bean)
		})
	}
	for _, v := range downloadList {
		doc, _ := goquery.NewDocumentFromResponse(getResultHtml(v.href))
		doc.Find("source").Each(func(i int, contentSelection *goquery.Selection) {
			download_url, _ := contentSelection.Attr("src")
			one_file := "wget " + download_url + "  -O \"./" + v.dirname + "/" + v.filename + ".mp4\"\n"
			bash_str += one_file
			fmt.Println(one_file)
		})
	}
	err := ioutil.WriteFile("./download.sh", []byte(bash_str), 0x777)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("写入脚本完成")
}
