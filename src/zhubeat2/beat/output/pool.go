package output

import (
	"errors"
	"fmt"
	"time"
	"zhubeat/lib/transport"
)

type ClientPool struct {
	total        int                // 票的总数。
	ticketClient chan *OutputClient // 票的容器。
	active       bool               // 票池是否已被激活。
}

// NewGoTickets 会新建一个Goroutine票池。
func NewClientPool(network string, hosts []string, timeout, ttl time.Duration, total int) (*ClientPool, error) {
	pool := &ClientPool{}
	if !pool.init(network, hosts, timeout, ttl, total) {
		errMsg :=
			fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return pool, nil
}

func (pool *ClientPool) init(network string, hosts []string, timeout, ttl time.Duration, total int) bool {

	if pool.active {
		return false
	}
	if total == 0 {
		return false
	}
	ch := make(chan *OutputClient, total)
	n := int(total)
	for i := 0; i < n; i++ {
		conn, err := transport.NewClient(network, hosts, timeout)
		if err != nil {
			return false
		}
		outClient := NewOutputClient(conn, ttl)
		err = outClient.Connect()
		if err != nil {
			// handle err
			return false
		}
		ch <- outClient
	}
	pool.ticketClient = ch
	pool.total = total
	pool.active = true
	return true
}

func (pool *ClientPool) Take() *OutputClient {
	return <-pool.ticketClient
}

func (pool *ClientPool) Return(client *OutputClient) {
	pool.ticketClient <- client
}

func (pool *ClientPool) Active() bool {
	return pool.active
}

func (pool *ClientPool) Total() int {
	return pool.total
}

func (pool *ClientPool) Remainder() int {
	return len(pool.ticketClient)
}

func (pool *ClientPool) Close() {
	for i := 0; i < pool.Total(); i++ {
		client := <-pool.ticketClient
		client.Close()
	}
}
