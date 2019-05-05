package main

import (
	"video-crawler/crawler-concurrent-show/engine"
	"video-crawler/crawler-concurrent-show/persist"
	"video-crawler/crawler-concurrent-show/scheduler"
	"video-crawler/crawler-concurrent-show/zhenai/parser"
)

func main() {
	//http://www.zhenai.com/zhenghun
	itemChan, err := persist.ItemSaver("dating_profile")
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		 : engine.Worker,
	}
	e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun", ParserFunc: engine.NewFuncParser(parser.ParseCityList, "ParseCityLists")})

}
