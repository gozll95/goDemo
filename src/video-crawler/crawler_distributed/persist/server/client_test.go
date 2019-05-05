package main

import (
	"fmt"
	"testing"
	"time"
	"video-crawler/crawler-concurrent-show/model"

	"video-crawler/crawler-concurrent-show/engine"
	"video-crawler/crawler_distributed/config"
	"video-crawler/crawler_distributed/rpcsupport"
)

func TestItemSaver(t *testing.T) {
	const host = fmt.Sprintf(":%d", config.ItemSaverPort)
	//start ItemSaverServer
	go serveRpc(host, "test1")

	time.Sleep(time.Second)
	//start ItemSaverClient
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}

	//call save
	item := engine.Item{
		Url:  "http://album.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: model.Profile{
			Name:     "HiSiri",
			Gender:   "女",
			Age:      28,
			Height:   163,
			Weight:   100,
			Income:   "3001-5000元",
			Marriage: "未婚",
			Hukou:    "内蒙古赤峰",
			Xingzuo:  "金牛座",
			House:    "自住",
			Car:      "未购车",
		},
	}
	result := ""
	err = client.Call(config.ItemSaverRpc, item, &result)
	if err != nil || result != "ok" {
		t.Errorf("result：%s %s", err, result)
	}
}
