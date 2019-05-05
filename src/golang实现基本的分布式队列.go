/*
http://blog.csdn.net/xcl168/article/details/44590457


 最近研究了下gowoker，这东西代码少而精，Golang真是很适合实现这类东西。
我去掉参数配置,JSON,Redis这些东西，用goworker的方式做了个最简单的实现。

  实现如下功能:
     1. worker向JobServer注册可执行的功能
     2. JobServer轮询,有job就执行,没有则继续轮询
     3. client向JobServer提出任务请求,并传入参数
     4. JobServer依请求丢给worker执行(可并发或串行执行)
     5. JobServer继续轮询

        我弄的这个代码很少，其中队列用数组代替,同时省掉了很多东西,
但保留了其goroutine与channel最基础的实现。
如果想看goworker的,可以参考下我这个,应当可以更快的弄明白goworker。


*/

package main

//分布式后台任务队列模拟(一)
//author: Xiong Chuan Liang
//date: 2015-3-24

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type workerFunc func(string, ...interface{}) error

type Workers struct {
	workers map[string]workerFunc
}

type OrdType int

const (
	PARALLEL = 1 << iota
	ORDER
)

type JobServer struct {
	Workers
	JobQueue []*WorkerClass
	interval time.Duration
	mt       sync.Mutex
	ord      OrdType
}

func NewJobServer() *JobServer {
	s := &JobServer{}
	s.workers = make(map[string]workerFunc, 0)
	return s
}

func (s *JobServer) RegisterWorkerClass(className string, f workerFunc) int {
	s.mt.Lock()
	defer s.mt.Unlock()
	if _, found := s.workers[className]; found {
		return 1
	}
	s.workers[className] = f
	return 0
}

type WorkerClass struct {
	ClassName string
	Args      []interface{}
}

func (s *JobServer) Enqueue(className string, args ...interface{}) bool {
	s.mt.Lock()
	w := &WorkerClass{className, args}
	s.JobQueue = append(s.JobQueue, w)
	s.mt.Unlock()
	return true
}

//poller
func (s *JobServer) poll(quit <-chan bool) <-chan *WorkerClass {
	jobs := make(chan *WorkerClass)

	go func() {
		defer close(jobs)
		for {
			switch {
			case s.JobQueue == nil:
				timeout := time.After(time.Second * 2)
				select {
				case <-quit:
					fmt.Println("[JobServer] [poll] quit")
					return
				case <-timeout:
					fmt.Println("[JobServer] [poll] polling")
				}
			default:

				s.mt.Lock()
				j := s.JobQueue[0]
				if len(s.JobQueue)-1 <= 0 {
					s.JobQueue = nil
				} else {
					s.JobQueue = s.JobQueue[1:len(s.JobQueue)]
				}
				s.mt.Unlock()

				select {
				case jobs <- j:
				case <-quit:
					fmt.Println("[JobServer] [poll] quit")
					return
				}

			}
		}
	}()
	return jobs
}

//worker
func (s *JobServer) work(id int, jobs <-chan *WorkerClass, monitor *sync.WaitGroup) {
	monitor.Add(1)

	f := func() {
		defer monitor.Done()
		for job := range jobs {
			if f, found := s.workers[job.ClassName]; found {
				s.run(f, job)
			} else {
				fmt.Println("[JobServer] [poll] ", job.ClassName, " not found")
			}
		}
	}

	switch s.ord {
	case ORDER:
		f()
	default:
		go f()
	}
}

func (s *JobServer) run(f workerFunc, w *WorkerClass) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[JobServer] [run] Panicking %s\n", fmt.Sprint(r))
		}
	}()

	f(w.ClassName, w.Args...)
}

func (s *JobServer) StartServer(interval time.Duration, ord OrdType) {

	s.interval = interval
	s.ord = ord

	quit := signals()
	// quit := make(chan bool)
	jobs := s.poll(quit)

	var monitor sync.WaitGroup

	switch s.ord {
	case ORDER: //顺序执行
		s.work(0, jobs, &monitor)
	default: //并发执行
		concurrency := runtime.NumCPU()
		for id := 0; id < concurrency; id++ {
			s.work(id, jobs, &monitor)
		}
	}

	monitor.Wait()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("分布式后台任务队列模拟(一)...")

	//Job Server
	js := NewJobServer()

	//模拟Worker端注册
	js.RegisterWorkerClass("mail", mailWorker)
	js.RegisterWorkerClass("log", sendLogWorker)
	js.RegisterWorkerClass("exception", paincWorker)

	//模拟客户端发送请求
	go func() {
		time.Sleep(time.Second * 2)
		js.Enqueue("mail", "xcl_168@aliyun.com", "sub", "body")

		js.Enqueue("test_notfound", "aaaaaaaaaaaaaaaaaaa")
		js.Enqueue("log", "x.log", "c.log", "l.log")

		//测试jobserver.PARALLEL/ORDER
		//for j := 0; j < 100; j++ {
		//	js.Enqueue("mail", strconv.Itoa(j))
		//}

		time.Sleep(time.Second)
		js.Enqueue("exception", "try{}exception{}")

		time.Sleep(time.Second * 5)
		js.Enqueue("mail", "xcl_168@aliyun.com2", "sub2", "body2")
	}()

	//启动服务，开始轮询
	// StartServer(轮询间隔,执行方式(并发/顺序))
	js.StartServer(time.Second*3, ORDER) //PARALLEL
}

func mailWorker(queue string, args ...interface{}) error {
	fmt.Println("......mail() begin......")
	for _, arg := range args {
		fmt.Println("   args:", arg)
	}
	fmt.Println("......mail() end......")
	return nil
}

func sendLogWorker(queue string, args ...interface{}) error {
	fmt.Println("......sendLog() begin......")
	for _, arg := range args {
		fmt.Println("   args:", arg)
	}
	fmt.Println("......sendLog() end......")
	return nil
}

func paincWorker(queue string, args ...interface{}) error {
	fmt.Println("......painc() begin......")
	panic("\n    test exception........................ \n")
	fmt.Println("......painc() end......")
	return nil
}

//有待商量
// func (s *JobServer) poll(quit <-chan bool) <-chan *WorkerClass {

func signals() <-chan bool {
	quit := make(chan bool)
	go func() {
		os.Stdin.Read(make([]byte, 1))
		quit <- true
	}()
	return quit
}

/*
flower@:~/workspace/learngo/src/myGoNotes$ go run golang实现基本的分布式队列.go
分布式后台任务队列模拟(一)...
[JobServer] [poll] polling
......mail() begin......
   args: xcl_168@aliyun.com
   args: sub
   args: body
......mail() end......
[JobServer] [poll]  test_notfound  not found
......sendLog() begin......
   args: x.log
   args: c.log
   args: l.log
......sendLog() end......
[JobServer] [poll] polling
......painc() begin......
[JobServer] [run] Panicking
    test exception........................


[JobServer] [poll] polling
[JobServer] [poll] polling
......mail() begin......
   args: xcl_168@aliyun.com2
   args: sub2
   args: body2
......mail() end......

[JobServer] [poll] polling
[JobServer] [poll] polling
[JobServer] [poll] polling
[JobServer] [poll] polling
[JobServer] [poll] polling

*/

//---->
// 上游函数里goroutine返回channel,下游range channel
