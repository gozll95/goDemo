package persist

import (
	"context"
	"errors"
	"log"

	"video-crawler/crawler-concurrent-show/engine"

	"gopkg.in/olivere/elastic.v5"
)

func ItemSaver(index string) (chan engine.Item, error) {
	client, err := elastic.NewClient(
		//Must turn off sniff in docker
		elastic.SetSniff(false),
	)
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
			err := Save(client, index, item)
			if err != nil {
				log.Printf("Item saver :error saving item %v:%v", item, err)
			}

		}
	}()
	return out, nil
}

func Save(client *elastic.Client, index string, item engine.Item) error {

	if item.Type == "" {
		return errors.New("must suplly type")
	}
	indexService := client.Index().Index(index).Type(item.Type).BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err := indexService.Do(context.Background())
	if err != nil {
		return err
	}
	//fmt.Printf("%+v", resp)
	return nil
}
