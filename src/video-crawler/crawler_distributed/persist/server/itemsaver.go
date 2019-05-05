package main

import (
	"flag"
	"fmt"
	"video-crawler/crawler_distributed/config"
	"video-crawler/crawler_distributed/persist"
	"video-crawler/crawler_distributed/rpcsupport"

	"gopkg.in/olivere/elastic.v5"
)

var port = flag.Int("port", 0,
	"the port for me to listen on")

func main() {
	flag.Parse()
	fetcher.SetVerboseLogging()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}

	serveRpc(fmt.Sprintf(":%d", *port), config.ElasticIndex)
}

// func main() {
// 	serveRpc(fmt.Sprintf(":%d", config.ItemSaverPort), config.ElasticIndex)
// }

func serveRpc(host, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}
	return rpcsupport.ServerRpc(host, &persist.ItemSaverService{
		Client: client,
		Index:  index,
	})
}
