package main

import (
	"video-crawler/crawler-single/engine"
	"video-crawler/crawler-single/zhenai/parser"
)

func main() {
	engine.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})

}
