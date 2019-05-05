package main

import (
	"fmt"
	"video-crawler/crawler-concurrent-show/engine"
	"video-crawler/crawler-concurrent-show/scheduler"
	"video-crawler/crawler-concurrent-show/zhenai/parser"
	"video-crawler/crawler_distributed/config"
	itemsaver "video-crawler/crawler_distributed/persist/client"
	worker "video-crawler/crawler_distributed/worker/client"
	"video-crawler/crawler_distributed/rpcsupport"
	"net/rpc"
	"flag"
	"strings"
)


var(
	itemSaverHost=flag.String(
		"itemsaver_host","","itemsaver host" )
	workerHost=flag.String(
		"worker_host","","worker hosts (comma separated)"
	)
)
func main() {
	flag.Parse()
	//http://www.zhenai.com/zhenghun
	itemChan, err := itemsaver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}
	


	pool:=createClientPool(
		strings.Split(*workerHost,",")
	)
	processor, err := worker.CreateProcessor(pool)
	if err != nil {
		panic(err)
	}

	e := engine. {
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: processor, // 现在是100个worker共用这个processor,怎么让这100个人去连不同的processor
		//RequestProcessorCreator: xxx, 可以每个人来带不同的参数来创建不同的processor
	}
	e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun",
		ParserFunc: engine.NewFuncParser(parser.ParseCityList, config.ParseCityList)})

}


func createClientPool(hosts []string )chan *rpc.Client{
   var clients []*rpc.Client
   for _,h:=range hosts{
	   client,err:=rpcsupport.NewClient(h)
	   if err==nil{
		   clients=append(clients,client)
		    log.Printf("Connected to %s",h)
	   }else{
		   log.Printf("error connecting to %s:%v",h,err)
	   }
   }
   out:=make(chan *rpc.Client)
   go func(){
		for{
			// 开始分发 可以随机分发 也可以轮流分发
			for _,client:=range clients{
				out<-client
			}
		}
   }()
   return out
}