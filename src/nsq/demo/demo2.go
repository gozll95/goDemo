package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	producer()
	go consumer() // simulate two consumer
	consumer2()
}
func producer() {
	producer, err := nsq.NewProducer("192.168.33.10:4150", nsq.NewConfig())
	defer producer.Stop()
	if err != nil {
		fmt.Println(err.Error())
	}
	for i := 0; i < 10000; i++ {
		// create a topic named testTopic
		err = producer.Publish("testTopic", []byte("testing...."+strconv.Itoa(i)))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
func consumer() {
	// creat a consumer with channel 'channelTestOne'
	consumer, err := nsq.NewConsumer("testTopic", "channelTestOne", nsq.NewConfig())
	if err != nil {
		fmt.Println(err.Error())
	}
	handler := new(NSQMessageHandler)
	handler.msgchan = make(chan *nsq.Message, 1024)
	consumer.AddHandler(nsq.HandlerFunc(handler.HandleMessage))
	err = consumer.ConnectToNSQLookupd("192.168.33.10:4161")
	if err != nil {
		fmt.Println(err.Error())
	}
	handler.Process()
}
func consumer2() {
	// creat another consumer with channel 'channelTestTwo'
	consumer, err := nsq.NewConsumer("testTopic", "channelTestTwo", nsq.NewConfig())
	if err != nil {
		fmt.Println(err.Error())
	}
	handler := new(NSQMessageHandler)
	handler.msgchan = make(chan *nsq.Message, 1024)
	consumer.AddHandler(nsq.HandlerFunc(handler.HandleMessage))
	err = consumer.ConnectToNSQLookupd("192.168.33.10:4161")
	if err != nil {
		fmt.Println(err.Error())
	}
	handler.Process()
}

type NSQMessageHandler struct {
	msgchan chan *nsq.Message
	stop    bool
}

func (m *NSQMessageHandler) HandleMessage(message *nsq.Message) error {
	if !m.stop {
		m.msgchan <- message
	}
	return nil
}
func (m *NSQMessageHandler) Process() {
	m.stop = false
	for {
		select {
		case message := <-m.msgchan:
			fmt.Println(string(message.Body))
		case <-time.After(time.Second):
			if m.stop {
				close(m.msgchan)
				return
			}
		}
	}
}
