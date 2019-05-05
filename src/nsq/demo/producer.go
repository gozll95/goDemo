package main

import (
	"github.com/nsqio/go-nsq"
)

var producer *nsq.Producer

func main() {
	nsqd := "127.0.0.1:4150"
	producer, err := nsq.NewProducer(nsqd, nsq.NewConfig())
	producer.Publish("testtt", []byte("nihao"))
	if err != nil {
		panic(err)
	}
}
