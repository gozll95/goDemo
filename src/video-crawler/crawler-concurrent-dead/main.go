package main

import (
	"video-crawler/crawler-concurrent-dead/engine"
	"video-crawler/crawler-concurrent-dead/scheduler"
	"video-crawler/crawler-concurrent-dead/zhenai/parser"
)

func main() {
	// engine.SimpleEngine{}.Run(engine.Request{
	// 	Url:        "http://www.zhenai.com/zhenghun",
	// 	ParserFunc: parser.ParseCityList,
	// })

	e := &engine.ConcurrentEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 10,
	}
	e.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})

}

/*
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun
2018/05/19 15:10:22 Got item: City:阿坝
2018/05/19 15:10:22 Got item: City:阿克苏
2018/05/19 15:10:22 Got item: City:阿拉善盟
2018/05/19 15:10:22 Got item: City:阿勒泰
2018/05/19 15:10:22 Got item: City:阿里
2018/05/19 15:10:22 Got item: City:安徽
2018/05/19 15:10:22 Got item: City:安康
2018/05/19 15:10:22 Got item: City:安庆
2018/05/19 15:10:22 Got item: City:鞍山
2018/05/19 15:10:22 Got item: City:安顺
2018/05/19 15:10:22 Got item: City:安阳
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/anshun
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/ali
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/aba
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/alashanmeng
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/ankang
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/anshan
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/aletai
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/anqing
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/anhui
2018/05/19 15:10:22 Fetching  http://www.zhenai.com/zhenghun/akesu
*/

//会一直卡着

/*
因为
scheduler 送 request 给 worker
worker 收 request 发 result 给 engine
engine 收 result 发 request 给 scheduler
scheduler 送 request 给 worker

以上是一个循环等待

会造成卡死
*/
