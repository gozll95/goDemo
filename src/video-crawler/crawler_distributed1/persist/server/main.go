package main

import (
	"video-crawler/crawler_distributed/persist"
	"video-crawler/crawler_distributed/rpcsupport"

	"gopkg.in/olivere/elastic.v5"
)

func main() {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	rpcsupport.ServeRpc(":1234",
		persist.ItemSaverService{
			Client: client,
			Index:  "dating_profile",
		})
}
