package scheduler

import (
	"video-crawler/crawler-concurrent/engine"
)

type QueuedScheduler struct {
	requestChan chan engine.Request
	workerChan  chan chan engine.Request //实际上是 chan worker,worker对外的接口是 chan engine.Request
}

func (s *QueuedScheduler) Submit(r engine.Request) {
	s.requestChan <- r
}

func (s *QueuedScheduler) WorkerReady(w chan engine.Request) {
	s.workerChan <- w
}

func (s *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request) // 希望每个worker拥有一个自己的channel request
}

// func (s *QueuedScheduler) ConfigureMasterWorkerChan(r chan engine.Request) {

// }

func (s *QueuedScheduler) Run() {
	s.workerChan = make(chan chan engine.Request)
	s.requestChan = make(chan engine.Request)
	go func() {
		var (
			requestQ []engine.Request
			workerQ  []chan engine.Request
		)
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			select {
			case r := <-s.requestChan:
				// send r to a ?worker
				// 我们可以把它放入队列
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				// send ?next_request to w
				// 我们可以把它放入队列
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest: // 如果并没有activeWorker,那么这一行是nil。这一行永远不会被selected
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}
