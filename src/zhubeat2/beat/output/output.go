package output

import (
	"fmt"
	"time"
	"zhubeat/beat"
	"zhubeat/beat/queue"
	lib_queue "zhubeat/lib/queue"
)

type OutPutManager struct {
	hosts       []string
	timeout     time.Duration
	ttl         time.Duration
	clientPool  *ClientPool
	concurrency int // 并发量
	queue       *lib_queue.Queue
	quit        chan struct{}
}

func NewOutPutManager(network string, hosts []string, timeout, ttl time.Duration, concurrency int, queue *lib_queue.Queue) (*OutPutManager, error) {
	clientPool, err := NewClientPool(network, hosts, timeout, ttl, concurrency)
	if err != nil {
		panic(err)
	}

	return &OutPutManager{
		hosts:       hosts,
		ttl:         ttl,
		timeout:     timeout,
		concurrency: concurrency,
		clientPool:  clientPool,
		queue:       queue,
		quit:        make(chan struct{}),
	}, nil
}

func (o *OutPutManager) Run() {
	defer func() {

	}()
	select {
	case <-o.quit:
		return
	default:
	}
	go o.start()

}

func (o *OutPutManager) start() {
	for {
		select {
		case <-o.quit:
			return
		default:
			if queue.ToOutputQueue.Len() > 0 {
				client := o.clientPool.Take()
				go func(client *OutputClient) {
					defer o.clientPool.Return(client)
					value, err := queue.ToOutputQueue.Pop()
					if err != nil {
						return
					}

					var jobs []beat.Job
					var ok bool

					if jobs, ok = value.([]beat.Job); !ok {
						fmt.Println("xxx")
						return
					}

					// handle err
					// conn, err := transport.NewClient("tcp", o.hosts, o.timeout)
					// outClient := NewOutputClient(conn, o.ttl)
					err = client.Publish(jobs)
					if err != nil {
						panic(err)
					}
				}(client)
			}
		}
	}
}

func (o *OutPutManager) Stop() {
	close(o.quit)
	o.clientPool.Close()
}
