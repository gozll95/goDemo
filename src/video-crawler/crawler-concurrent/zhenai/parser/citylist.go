package parser

import (
	"regexp"
	"video-crawler/crawler-concurrent/engine"
)

const cityLisRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

/**
  正则匹配所有 url
*/
func ParseCityList(contents []byte) engine.ParseResult {
	re := regexp.MustCompile(cityLisRe)
	matches := re.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}
	limit := 10

	for _, m := range matches {
		result.Items = append(result.Items, "City:"+string(m[2]))
		result.Requests = append(result.Requests,
			engine.Request{
				Url:        string(m[1]),
				ParserFunc: ParseCity})
		// fmt.Printf("City: %s Url : %s\n",
		// 	m[2], m[1])
		limit--
		if limit < 0 {
			break
		}
	}

	// fmt.Printf("Matches found: %d\n",
	// 	len(matches))
	return result
}
