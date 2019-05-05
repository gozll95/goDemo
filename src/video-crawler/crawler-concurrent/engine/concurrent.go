package engine

import "log"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type Scheduler interface {
	// simple scheduler
	// Submit(Request)
	// ConfigureMasterWorkerChan(chan Request)

	// 这是对于queue.go
	// Submit(Request)
	// ConfigureMasterWorkerChan(chan Request)
	// WorkerReady(chan Request)
	// Run()

	// 重构之后
	Submit(Request)
	WorkerChan() chan Request
	//WorkerReady(chan Request)
	Run()

	ReadyNotifier //使用interface组合的方式
}

// 重构添加
type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult)
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		// 我不知道每个worker一个channel还是所有worker共用一个channel,这个问题Scheduler知道,所以我问它要
		createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	itemCount := 0
	for {
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item#%d: %v", itemCount, item)
			itemCount++
		}
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

//func createWorker(in chan Request, out chan ParseResult, s Scheduler) {
func createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			// tell scheduler i`m ready
			ready.WorkerReady(in)
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
