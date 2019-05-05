package main

import (
	"fmt"
	"net/url"
)

func main() {
	values, err := url.ParseRequestURI("https://www.baidu.com/s?wd=%E6%90%9C%E7%B4%A2&rsv_spt=1&issp=1&f=8&rsv_bp=0&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_sug3=7&rsv_sug1=6")

	fmt.Println(values)
	// 会打印出https://www.baidu.com/s?wd=%E6%90%9C%E7%B4%A2&rsv_spt=1&issp=1&f=8&rsv_bp=0&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_sug3=7&rsv_sug1=6

	if err != nil {
		fmt.Println(err)
	}

	urlParam := values.RawQuery
	fmt.Println(urlParam) // 打印出wd=%E6%90%9C%E7%B4%A2&rsv_spt=1&issp=1&f=8&rsv_bp=0&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_sug3=7&rsv_sug1=6

	// ParseQuery传入的必须是参数，也就是url里边的RawQuery的值 就是url?之后的path
	fmt.Println(url.ParseQuery(urlParam))
	//会打印出map[tn:[baiduhome_pg] rsv_enter:[1] rsv_sug3:[7] rsv_sug1:[6] f:[8] rsv_bp:[0] ie:[utf-8] rsv_idx:[2] wd:[搜索] rsv_spt:[1] issp:[1]] <nil>

	//这里url.Query()直接就解析成map了，url.ParseQuery()反而多了一步，果断用这个方法
	urlValue := values.Query() // 和下面的c变量类型相同都为url.Values类型，有相同的属性方法
	fmt.Println(urlValue)
	// 会打印出map[rsv_bp:[0] ie:[utf-8] tn:[baiduhome_pg] rsv_enter:[1] rsv_sug3:[7] rsv_sug1:[6] f:[8] rsv_spt:[1] issp:[1] rsv_idx:[2] wd:[搜索]]

	//val := url.Values{}
	c := url.Values{"method": {"get"}, "id": {"1"}}
	fmt.Println(c.Encode()) // 打印出id=1&method=get
	c.Get("method")         // 获取到method的值为get

	c.Set("method", "post") // 修改method的值为post

	c.Del("method") // 删除method元素

	c.Add("new", "hi") // 添加新的元素new:hi
}
