package scheduler

import (
	"video-crawler/crawler-concurrent-queue-scheduler/engine"
)

type SimpleScheduler struct {
	workerChan chan engine.Request
}

func (s *SimpleScheduler) Submit(r engine.Request) {
	// send request down to worker chan
	//s.workerChan <- r //这个会因为循环等待而卡死
	go func() { // 所以改造成goroutine这种模式
		s.workerChan <- r
	}()
}

func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request) {
	s.workerChan = c
}
