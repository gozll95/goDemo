package persist

import (
	"context"
	"encoding/json"
	"testing"

	"video-crawler/crawler-concurrent-show/engine"
	"video-crawler/crawler-concurrent-show/model"

	"gopkg.in/olivere/elastic.v5"
)

func TestSaver(t *testing.T) {
	expected := engine.Item{
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
	const index = "dating_test"
	client, err := elastic.NewClient(
		//Must turn off sniff in docker
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	//save expected
	err = Save(client, index, expected)
	if err != nil {
		panic(err)
	}

	//fetch expected
	resp, err := client.Get().Index(index).Type(expected.Type).Id(expected.Id).Do(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("%+v,%s", resp, resp.Source)
	var actual engine.Item
	err = json.Unmarshal([]byte(*resp.Source), &actual)
	if err != nil {
		panic(err)
	}
	actualProfile, err := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile
	//verify result
	if actual != expected {
		t.Errorf("got %v;expected %v", actual, expected)
	}
}
