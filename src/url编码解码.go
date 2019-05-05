package main

import (
	"fmt"
	"net/url"
)

func main() {
	urltest := "http://www.baidu.com/s?wd=自由度"
	fmt.Println(urltest)
	encodeurl := url.QueryEscape(urltest)
	fmt.Println(encodeurl)
	decodeurl, err := url.QueryUnescape(encodeurl)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(decodeurl)

}

/*
http://www.baidu.com/s?wd=自由度
http%3A%2F%2Fwww.baidu.com%2Fs%3Fwd%3D%E8%87%AA%E7%94%B1%E5%BA%A6
http://www.baidu.com/s?wd=自由度
*/
