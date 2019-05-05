package output

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"zhubeat/lib/transport"
)

type WorkerPool struct {
	total        int
	ticketClient chan *OutputClient
	closed       uint32
	rwlock       sync.RWMutex
}

func NewWorkerPool(network string, hosts []string, timeout, ttl time.Duration, total int) (*WorkerPool, error) {
	pool := &WorkerPool{}
	if !pool.init(network, hosts, timeout, ttl, total) {
		errMsg :=
			fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return pool, nil
}

func (pool *WorkerPool) init(network string, hosts []string, timeout, ttl time.Duration, total int) bool {

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
		err = outClient.RandomConnect()
		if err != nil {
			// handle err
			return false
		}
		ch <- outClient
	}
	pool.ticketClient = ch
	pool.total = total
	return true
}

func (pool *WorkerPool) Take() *OutputClient {
	if pool.Closed() {
		return nil
	}

	pool.rwlock.RLock()
	defer pool.rwlock.RUnlock()
	return <-pool.ticketClient
}

func (pool *WorkerPool) Return(client *OutputClient) {
	if pool.Closed() {
		return
	}
	pool.rwlock.RLock()
	defer pool.rwlock.RUnlock()

	pool.ticketClient <- client
}

func (pool *WorkerPool) Total() int {
	return pool.total
}

func (pool *WorkerPool) Remainder() int {
	return len(pool.ticketClient)
}

func (pool *WorkerPool) Close() bool {
	if !atomic.CompareAndSwapUint32(&pool.closed, 0, 1) {
		return false
	}
	pool.rwlock.Lock()
	defer pool.rwlock.Unlock()
	close(pool.ticketClient)
	for client := range pool.ticketClient {
		client.Close()
	}
	return true
}

func (pool *WorkerPool) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}
