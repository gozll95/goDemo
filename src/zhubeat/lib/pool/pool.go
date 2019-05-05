package pool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type PoolClient interface {
	Close() error
}
type Pool struct {
	total         int
	clients       chan PoolClient
	closed        uint32
	rwlock        sync.RWMutex
	newClientFunc func() (PoolClient, error)
}

func NewPool(total int, newClientFunc func() (PoolClient, error)) (*Pool, error) {
	pool := &Pool{}
	pool.newClientFunc = newClientFunc
	if !pool.init(total) {
		errMsg :=
			fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return pool, nil
}

func (pool *Pool) init(total int) bool {
	if total == 0 {
		return false
	}
	ch := make(chan PoolClient, total)
	n := int(total)
	for i := 0; i < n; i++ {
		client, err := pool.newClientFunc()
		if err != nil {
			return false
		}
		ch <- client
	}
	pool.clients = ch
	pool.total = total
	return true
}

func (pool *Pool) TakeClient() PoolClient {
	if pool.Closed() {
		return nil
	}

	pool.rwlock.RLock()
	defer pool.rwlock.RUnlock()
	return <-pool.clients
}

func (pool *Pool) ReturnClient(client PoolClient) {
	if pool.Closed() {
		return
	}
	pool.rwlock.RLock()
	defer pool.rwlock.RUnlock()

	pool.clients <- client
}

func (pool *Pool) Total() int {
	return pool.total
}

func (pool *Pool) Remainder() int {
	return len(pool.clients)
}

func (pool *Pool) Close() bool {
	if !atomic.CompareAndSwapUint32(&pool.closed, 0, 1) {
		return false
	}
	pool.rwlock.Lock()
	defer pool.rwlock.Unlock()
	close(pool.clients)
	for client := range pool.clients {
		client.Close()
	}
	return true
}

func (pool *Pool) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}
