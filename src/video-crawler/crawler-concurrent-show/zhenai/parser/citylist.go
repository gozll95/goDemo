package parser

import (
	"regexp"
	"video-crawler/crawler-concurrent-show/engine"
)

const citylistre = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParseCityList(contents []byte, _ string) engine.ParseResult {
	re := regexp.MustCompile(citylistre)
	match := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, m := range match {
		url := string(m[1])
		//result.Items = append(result.Items, "City "+string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:        url,
			ParserFunc: engine.NewFuncParser(ParseCity, "ParseCity"),
		})
		//fmt.Printf("%s\n %s \n", m[1], m[2])
	}
	return result
}
