package client

import (
	"log"
	"video-crawler/crawler-concurrent-show/engine"
	"video-crawler/crawler_distributed/config"
	"video-crawler/crawler_distributed/rpcsupport"
)

func ItemSaver(host string) (chan engine.Item, error) {
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			itemCount++
			item := <-out
			log.Printf("ItemSaver got Item: #%d:%v", itemCount, item)

			// Call RPC to save item
			result := ""
			err = client.Call(config.ItemSaverRpc, item, &result)
			if err != nil {
				log.Printf("Item saver :error saving item %v:%v", item, err)
			}

		}
	}()
	return out, nil
}
