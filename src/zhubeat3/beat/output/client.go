package output

import (
	"fmt"
	"log"
	"time"
	"zhubeat/beat/job"
	"zhubeat/lib/transport"
)

type ClientArgs struct {
	network   string
	addresses []string // ["127.0.0.1:7070","127.0.0.1:7070","127.0.0.1:7070"]
	timeout   time.Duration
	ttl       time.Duration
	// retry
}

func NewClientArgs(network string, addresses []string, timeout, ttl time.Duration) *ClientArgs {
	return &ClientArgs{
		network:   network,
		addresses: addresses,
		timeout:   timeout,
		ttl:       ttl,
	}
}

type OutputClient struct {
	*transport.Client
	ttl time.Duration
	// retry
}

func NewOutputClient(client *transport.Client, ttl time.Duration) *OutputClient {
	c := &OutputClient{
		Client: client,
		ttl:    ttl,
	}

	return c
}

// func (c *OutputClient) Connect() error {
// 	return c.Client.RandomConnect()
// }

// func (c *OutputClient) Close() error {
// 	return c.Client.Close()
// }

func (c *OutputClient) reconnect() error {
	if err := c.Client.Close(); err != nil {
		log.Println("error closing connection , reconnecting...", err)
	}
	return c.Client.RandomConnect()
}

// 每隔xx时间发一行
// 每次发一行
// retry机制
// retry失败 落地机制 // 可以忽略 tcp有重试机制
// failover???
// keepalive????
func (c *OutputClient) Publish(batch job.BatchJob) error {
	//	defer c.Close()

	// // set c.conn
	// err := c.Connect()
	// if err != nil {
	// 	return err
	// }

	for _, event := range batch {
		_, err := c.Write([]byte(string(event)))
		if err != nil {
			fmt.Println("store event to file or db ", string(event))
			fmt.Println(err)
			// event store to file/db
			err = c.reconnect()
			if err != nil {
				fmt.Println(err)
			}
		}
		time.Sleep(c.ttl)

	}
	return nil
}
