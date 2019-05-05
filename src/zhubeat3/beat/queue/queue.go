package queue

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"zhubeat/beat/job"
	"zhubeat/lib/queue"
	"zhubeat/lib/watcher"
)

var LogQueue *JobQueue

type JobQueue struct {
	queue     *queue.Queue
	lock      sync.RWMutex
	batchLen  int
	threshold int
	ttl       time.Duration

	// watched 代表监听器的关闭状态：0-未监听；1-已监听。
	watched   uint32
	watcher   *watcher.Watcher
	watchOnce sync.Once
}

func NewJobQueue(batchLen, threshold int, ttl time.Duration, queue *queue.Queue) *JobQueue {
	return &JobQueue{
		queue:     queue,
		threshold: threshold,
		ttl:       ttl,
		batchLen:  batchLen,
	}
}

func (q *JobQueue) Push(job job.Job) {
	defer q.lock.Unlock()
	q.lock.Lock()

	q.queue.Push(job)
}

func (q *JobQueue) Pop() (job.Job, error) {
	defer q.lock.Unlock()
	q.lock.Lock()
	return q.pop()
}

func (q *JobQueue) pop() (job.Job, error) {
	value, err := q.queue.Pop()
	if err != nil {
		return "", err
	}
	return value.(job.Job), nil
}

func (q *JobQueue) PopBatch() (batch job.BatchJob) {
	defer q.lock.Unlock()
	q.lock.Lock()

	for i := 0; i < q.batchLen; i++ {
		job, err := q.pop()
		if err != nil {
			return
		}
		batch = append(batch, job)
	}
	return
}

func (q *JobQueue) Len() int {
	defer q.lock.RUnlock()
	q.lock.RLock()

	return q.queue.Len()
}

func (q *JobQueue) Dump() (batch job.BatchJob) {
	defer q.lock.RUnlock()
	q.lock.RLock()

	return q.dump()
}

func (q *JobQueue) dump() (batch job.BatchJob) {
	values := q.queue.Dump()
	for _, v := range values {
		batch = append(batch, v.(job.Job))
	}
	return batch
}

func (q *JobQueue) Reset() {
	defer q.lock.Unlock()
	q.lock.Lock()

	q.queue.Reset()
}

func (q *JobQueue) DumpAndReset() job.BatchJob {
	defer q.lock.Unlock()
	q.lock.Lock()

	batch := q.dump()
	q.queue.Reset()

	return batch
}

func (q *JobQueue) IsOverThreshold() bool {
	if q.Len() > q.threshold {
		return true
	}
	return false
}

func (q *JobQueue) OverThresholdProccess() error {
	batch := q.DumpAndReset()
	fmt.Println("watcher store batch: ", batch)
	// store
	// 情空
	return nil
}

func (q *JobQueue) RunWatcher() {
	q.watchOnce.Do(func() {
		q.watcher = watcher.NewWatcher(q.ttl, q)
		q.watcher.Start()
		atomic.StoreUint32(&q.watched, 1)
	})

}

func (q *JobQueue) CloseWatcher() {
	if !atomic.CompareAndSwapUint32(&q.watched, 1, 0) {
		return
	}

	q.watcher.Close()
}

func (q *JobQueue) IsWatched() uint32 {
	return atomic.LoadUint32(&q.watched)
}
