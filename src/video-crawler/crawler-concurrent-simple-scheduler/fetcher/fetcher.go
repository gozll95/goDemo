package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}
	//*中文gbk -golang utf8
	// 如何发现文本是gbk呢? 1.可以使用meta charset="gbk" 2.使用
	//e := determineEncoding(resp.Body)
	//utf8Reader := transform.NewReader(resp.Body, e.NewDecoder())
	//return ioutil.ReadAll(utf8Reader)

	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// 会有中文乱码,中文gbk,golang是utf-8
// v1
// func determineEncoding(r io.Reader) encoding.Encoding {
// 	//这里很聪明,之前直接读io.Reader会导致下一次不能读,这里先将io.Reader
// 	//使用bufio缓存起来,然后使用Peek来返回一个Slice ???????
// 	// ??以上说法似乎不正确
// 	// 详情见:https://www.jianshu.com/p/a48298e4e2e2
// 	bytes, err := bufio.NewReader(r).Peek(1024)
// 	if err != nil {
// 		log.Printf("Fetcher error: %v", err)
// 		return unicode.UTF8 // 返回默认:utf-8
// 	}

// 	e, _, _ := charset.DetermineEncoding(bytes, "")
// 	return e
// }

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	//这里很聪明,之前直接读io.Reader会导致下一次不能读,这里先将io.Reader
	//使用bufio缓存起来,然后使用Peek来返回一个Slice ???????
	// ??以上说法似乎不正确
	// 详情见:https://www.jianshu.com/p/a48298e4e2e2
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8 // 返回默认:utf-8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
