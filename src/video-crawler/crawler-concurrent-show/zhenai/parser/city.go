package parser

import (
	"regexp"
	"video-crawler/crawler-concurrent-show/engine"
)

var (
	cityre    = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	cityUrlre = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/shanghai/[^"]+)">`)
)

func ParseCity(contents []byte, _ string) engine.ParseResult {
	amtch := cityre.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, m := range amtch {
		url := string(m[1])
		name := string(m[2])
		result.Requests = append(result.Requests, engine.Request{
			Url: url,
			//ParserFunc: ProfileParser(name),
			ParserFunc: NewProfileParser(name),
		})
		//fmt.Printf("%s\n %s \n", m[1], m[2])
	}
	matchs := cityUrlre.FindAllSubmatch(contents, -1)
	for _, m := range matchs {
		result.Requests = append(result.Requests, engine.Request{
			Url:        string(m[1]),
			ParserFunc: engine.NewFuncParser(ParseCity, "ParseCity"),
		})
	}
	return result
}
