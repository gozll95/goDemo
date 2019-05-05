package worker

import (
	"context"
	"fmt"
	"log"
	"zhubeat/beat"
	"zhubeat/beat/queue"
	"zhubeat/lib/manager"

	"time"

	"github.com/satori/go.uuid"
)

type Worker struct {
	Name       string
	WorkerPool chan chan manager.Job
	JobChannel chan manager.Job
	ctx        context.Context
}

func NewWorker(workerPool chan chan manager.Job, ctx context.Context) manager.Worker {
	name, _ := uuid.NewV4()
	return &Worker{
		Name:       name.String(),
		WorkerPool: workerPool,
		JobChannel: make(chan manager.Job),
		ctx:        ctx,
	}
}

func (w *Worker) Start() {
	ticker := time.Tick(1 * time.Second)

	var jobBatch []beat.Job
	defer func() {
		log.Println("closed worker:", w.Name)
	}()

	select {
	case <-w.ctx.Done():
		fmt.Println("receive done")
		return
	default:
	}
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println(w.Name, "receive done")
			return
		case w.WorkerPool <- w.JobChannel:
		case job := <-w.JobChannel:
			fmt.Println(w.Name)
			process(job)
			jobBatch = append(jobBatch, job.(beat.Job))
			if len(jobBatch) == 3 {
				send(jobBatch)
				jobBatch = []beat.Job{}
			}
		case <-ticker:
			// 这里其实 应该加锁
			if len(jobBatch) > 0 {
				send(jobBatch)
				jobBatch = []beat.Job{}
			}
		}
	}
}

type BatchJobs []manager.Job

func process(job manager.Job) {
	j := job.(beat.Job)
	fmt.Println(string(j))
}

func send(job []beat.Job) {
	fmt.Println("send length:", len(job))
	queue.ToOutputQueue.Push(job)
}
