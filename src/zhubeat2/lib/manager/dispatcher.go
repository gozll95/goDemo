package manager

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"zhubeat/lib/status"
)

//https://blog.lab99.org/post/golang-2017-10-27-video-how-to-correctly-use-package-context.html

type Dispatcher struct {
	name string

	// Job
	jobQueue chan Job

	// woker
	maxWorkers   int
	workerPool   chan chan Job
	wokerFactory WokerFactory

	//
	timeout    time.Duration
	ctx        context.Context
	cancelFunc context.CancelFunc

	// wait group
	wg sync.WaitGroup

	// status
	*status.StatusScheduler

	// extra
	// metric
	// monitor

	// new worker func

	// real-time active worker // some monitor info
	//summary SchedSummary
}

func NewDispatcher(maxWorkers int, jobQueue chan Job, wokerFactory WokerFactory) (d *Dispatcher) {
	log.Println("New a load dispatcher...")
	pool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		workerPool:      pool,
		name:            "dispatcher",
		maxWorkers:      maxWorkers,
		jobQueue:        jobQueue,
		wg:              sync.WaitGroup{},
		wokerFactory:    wokerFactory,
		StatusScheduler: status.NewStatusScheduler(),
	}
}

func (d *Dispatcher) Init() error {
	log.Println("dispatch init")
	var (
		oldStatus status.Status
		err       error
	)
	oldStatus, err = d.CheckAndSetStatus(status.SCHED_STATUS_INITIALIZING)
	if err != nil {
		panic(err)
		return err
	}

	d.ctx, d.cancelFunc = context.WithCancel(context.Background())
	// some checks

	defer func() {
		d.CheckAndSetStatusWithErr(oldStatus, status.SCHED_STATUS_INITIALIZED, err)
	}()

	log.Println("dispatch status: ", d.Status())
	return nil
}

func (d *Dispatcher) Start() error {
	var (
		oldStatus status.Status
		err       error
	)
	oldStatus, err = d.CheckAndSetStatus(status.SCHED_STATUS_STARTING)
	if err != nil {
		panic(err)
		return err
	}

	log.Println("dispatch status: ", d.Status())

	for i := 0; i < d.maxWorkers; i++ {
		worker := d.wokerFactory(d.workerPool, d.ctx)
		d.wg.Add(1)

		// 这里可以wrapper
		go func() {
			defer d.wg.Done()
			worker.Start()
		}()
	}
	go d.run()

	defer func() {
		d.CheckAndSetStatusWithErr(oldStatus, status.SCHED_STATUS_STARTED, err)
	}()

	return nil
}

func (d *Dispatcher) run() {
	for {
		select {
		case job := <-d.jobQueue:
			fmt.Println("调度者,接收到一个工作任务")
			//time.Sleep(time.Duration(500 * time.Millisecond))
			go func(job Job) {
				jobChannel := <-d.workerPool
				jobChannel <- job
			}(job)
		case <-d.ctx.Done():
			return
		default:

			//fmt.Println("ok!!")
		}

	}
}
func (d *Dispatcher) Status() status.Status {
	return d.StatusScheduler.Status()
}

func (d *Dispatcher) Close() {
	var (
		oldStatus status.Status
		err       error
	)
	oldStatus, err = d.CheckAndSetStatus(status.SCHED_STATUS_STOPPING)
	if err != nil {
		fmt.Println("err is ", err)
		return
	}

	d.cancelFunc()
	fmt.Println("cancelFunc")
	d.wg.Wait()

	// wait for queue chan?

	// wait for all worker close

	defer func() {
		d.CheckAndSetStatusWithErr(oldStatus, status.SCHED_STATUS_STOPPED, err)
	}()

}
