package parser

import (
	"io/ioutil"
	"testing"
	"video-crawler/crawler-concurrent-show/engine"
	"video-crawler/crawler-concurrent-show/model"
)

func TestParseProfile(t *testing.T) {
	contents, err := ioutil.ReadFile("profile_test_data.html")
	if err != nil {
		panic(err)
	}
	result := ParseProfile(contents, "HiSiri", "http://album.zhenai.com/u/108906739")
	if len(result.Items) != 1 {
		t.Error("Result should contain 1 element;but has %v", result.Items)
	}
	actula := result.Items[0]
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
	if actula != expected {
		t.Errorf("expected %v;but was %v", expected, actula)
	}

}
