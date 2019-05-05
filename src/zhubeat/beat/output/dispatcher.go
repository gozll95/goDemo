package output

import (
	"context"
	"fmt"
	"log"
	"time"
	"zhubeat/beat/job"
	"zhubeat/beat/queue"
	"zhubeat/lib/pool"
	"zhubeat/lib/status"
)

type Dispatcher struct {
	name string

	// Queue
	jobQueue *queue.JobQueue // pop出来的是batch job

	// woker
	maxWorkers int
	workerPool *pool.Pool

	// ticker
	ticker *time.Ticker

	ctx        context.Context
	cancelFunc context.CancelFunc

	// status
	*status.StatusScheduler
}

func NewDispatcher(maxWorkers int, jobQueue *queue.JobQueue, clientArgs *ClientArgs, ttl time.Duration) (d *Dispatcher) {
	log.Println("New a load dispatcher...")

	d = &Dispatcher{
		name:            "dispatcher",
		maxWorkers:      maxWorkers,
		jobQueue:        jobQueue,
		StatusScheduler: status.NewStatusScheduler(),
	}
	if ttl > 0 {
		d.ticker = time.NewTicker(ttl)
	}

	err := d.init(clientArgs)
	if err != nil {
		panic(err)
	}
	return
}

func (d *Dispatcher) init(clientArgs *ClientArgs) error {
	log.Println("dispatch init")
	var (
		oldStatus status.Status
		err       error
	)
	oldStatus, err = d.CheckAndSetStatus(status.SCHED_STATUS_INITIALIZING)
	defer func() {
		d.CheckAndSetStatusWithErr(oldStatus, status.SCHED_STATUS_INITIALIZED, err)
	}()
	if err != nil {
		panic(err)
		return err
	}

	// new pool
	connPool, err := NewConnPool(clientArgs.network, clientArgs.addresses)
	if err != nil {
		return err
	}

	// new client
	newclientFunc := func() (pool.PoolClient, error) {
		client := NewClient("", clientArgs.timeout, clientArgs.ttl, connPool)
		return client, nil
	}
	workerPool, err := pool.NewPool(d.maxWorkers, newclientFunc)
	if err != nil {
		panic(err)
	}
	d.workerPool = workerPool

	// init context
	d.ctx, d.cancelFunc = context.WithCancel(context.Background())

	log.Println("dispatch status: ", d.Status())
	return nil
}

func (d *Dispatcher) Start() error {
	var (
		oldStatus status.Status
		err       error
	)

	oldStatus, err = d.CheckAndSetStatus(status.SCHED_STATUS_STARTING)
	defer func() {
		d.CheckAndSetStatusWithErr(oldStatus, status.SCHED_STATUS_STARTED, err)
	}()

	if err != nil {
		panic(err)
		return err
	}

	log.Println("dispatch status: ", d.Status())

	go d.run()
	return nil
}

func (d *Dispatcher) run() {
	for {
		select {
		case <-d.ctx.Done():
			fmt.Println("receive done signal")
			return
		case <-d.ticker.C:
			jobBatch := d.jobQueue.PopBatch()
			if jobBatch.Len() > 0 {
				fmt.Println("get jobBatch: ", jobBatch)
				worker := d.workerPool.TakeClient().(*Client)
				go func(worker *Client, jobBatch job.BatchJob) {
					defer d.workerPool.ReturnClient(worker)

					// if !worker.Client.IsConnected() {
					// 	worker.Client.RandomConnect()
					// }

					worker.Publish(jobBatch)
				}(worker, jobBatch)
			}
		default:
		}
	}
}

func (d *Dispatcher) Status() status.Status {
	return d.StatusScheduler.Status()
}

func (d *Dispatcher) Close() {
	if d == nil {
		return
	}

	var (
		oldStatus status.Status
		err       error
	)
	oldStatus, err = d.CheckAndSetStatus(status.SCHED_STATUS_STOPPING)
	defer func() {
		d.CheckAndSetStatusWithErr(oldStatus, status.SCHED_STATUS_STOPPED, err)
	}()
	if err != nil {
		fmt.Println("err is ", err)
		return
	}

	d.cancelFunc()
	// for {
	// 	if d.workerPool.Remainder() == d.maxWorkers {
	// 		break
	// 	}
	// }

	d.workerPool.Close()

	//d.workerPool.Close()
}
