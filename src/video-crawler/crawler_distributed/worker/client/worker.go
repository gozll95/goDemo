package client

import (
	"net/rpc"
	"video-crawler/crawler-concurrent-show/engine"
	"video-crawler/crawler_distributed/config"
	"video-crawler/crawler_distributed/worker"
)

func CreateProcessor(
	clientChan chan *rpc.Client) engine.Processor {

	return func(
		req engine.Request) (
		engine.ParseResult, error) {

		sReq := worker.SerializeRequest(req)

		var sResult worker.ParseResult
		c := <-clientChan
		err := c.Call(config.CrawlServiceRpc,
			sReq, &sResult)

		if err != nil {
			return engine.ParseResult{}, err
		}
		return worker.DeserializeResult(sResult),
			nil
	}
}

// 第一版本
// func CreateProcessor() (engine.Processor, error) {
// 	client, err := rpcsupport.NewClient(fmt.Sprintf(":%d", config.WorkerPort0))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return func(r engine.Request) (engine.ParseResult, error) {
// 		sReq := worker.SerializeRequest(req)
// 		var sResult worker.ParseResult

// 		err = client.Call(config.CrawlServiceRpc, sReq, &sResult)
// 		if err != nil {
// 			return engine.ParseResult{}, err
// 		}
// 		return worker.DeserializeResult(sResult), nil

// 	}, nil
// }
