func main(){
    resp,err := login()
    if err != nil {
        print(err)
    }
    //for _,i := range resp.Cookies(){
    //    fmt.Println(i)
    //}
    url := "xxxxx"
    client := &http.Client{}
    req,err := http.NewRequest("POST",url,nil)
    if err != nil {
        print(err)
    }
    req.Header.Set("Cookie",resp.Cookies())  //数组报错
    req.Header.Set("Pragma","no-cache")
    req.Header.Set("Accept-Encoding","gzip, deflate, sdch")
    req.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
    req.Header.Set("Upgrade-Insecure-Requests","1")
    req.Header.Set("User-Agent","Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
    req.Header.Set("Accept","text/javascript, text/html, application/xml, text/xml, */*")
    req.Header.Set("Connection","keep-alive")
    req.Header.Set("Content-Type","application/x-www-form-urlencoded; charset=UTF-8")
    req.Header.Set("Cache-Control","no-cache")

    resp2, err := client.Do(req)
    if err != nil {
        print(err)
    }

    body, err := ioutil.ReadAll(resp2.Body)
    fmt.Println(string(body))
}

func login() (*http.Response, error) {
    url := "xxxxx"
    client := &http.Client{}
    req,err := http.NewRequest("POST",url,strings.NewReader("name=cjb"))
    if err != nil {
        print(err)
    }
    req.Header.Set("Pragma","no-cache")
    req.Header.Set("Accept-Encoding","gzip, deflate, sdch")
    req.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
    req.Header.Set("Upgrade-Insecure-Requests","1")
    req.Header.Set("User-Agent","Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
    req.Header.Set("Accept","text/javascript, text/html, application/xml, text/xml, */*")
    req.Header.Set("Connection","keep-alive")
    req.Header.Set("Content-Type","application/x-www-form-urlencoded; charset=UTF-8")
    req.Header.Set("Cache-Control","no-cache")
    //用户登陆设置...
    
    //请求,获取cookie
    resp, err := client.Do(req)

    //defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        print(err)
    }
    fmt.Println(string(body))
    return resp,nil