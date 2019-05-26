package main

import (
	"os"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type JobServer struct {
	Workers
	JobQueue []*WorkerClass
	interval time.Duration
	mt       sync.Mutex
	ord      OrdType
}

type Workers struct {
	workers map[string]workerFunc
}

type workerFunc func(string, ...interface{}) error

type WorkerClass struct {
	ClassName string
	Args      []interface{}
}

type OrdType int

const (
	PARALLEL = 1 << iota
	ORDER
)

func NewJobServer() *JobServer {
	s := &JobServer{}
	s.workers = make(map[string]workerFunc, 0)
	return s
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	js := NewJobServer()
	js.RegisterWorkerClass("mail", mailWorker)

	go func() {
		time.Sleep(time.Second * 2)
		js.Enqueue("mail", "xcl_168@aliyun.com", "sub", "body"

	}()
	js.StartServer(time.Second*3, ORDER) //PARALLEL
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

func mailWorker(queue string, args ...interface{}) error {
	fmt.Println("......mail() begin......")
	for _, arg := range args {
		fmt.Println("   args:", arg)
	}
	fmt.Println("......mail() end......")
	return nil
}

func (s *JobServer) Enqueue(className string, args ...interface{}) bool {
	s.mt.Lock()
	w := &WorkerClass{className, args}
	s.JobQueue = append(s.JobQueue, w)
	s.mt.Unlock()
	return true
}

func(s *JobServer)StartServer(interval time.Duration,ord OrdType){
	s.interval=interval
	s.ord=ord

	quit:=signals()

	jobs:=s.pool(quit)
}

func signals()  <- chan bool{
	quit:=make(chan bool)
	go func(){
		os.Stdin.Read(make([]byte),1)
		quit<-true
	}()
	return quit
}


func (s *JobServer)poll(quit <- chan bool)<-chan *WorkerClass{
	jobs:=make(chan *WorkerClass)

	go func(){
		defer close(jobs)
		for{
			switch{
			case s.JobQueue==nil:
				timeout := time.After(time.Second * 2)
				select{
				case<-quit:
				fmt.Println("[JobServer] [poll] quit")
					return	
				case <-timeout:
					fmt.Println("[JobServer] [poll] polling")
				}
			}
		}
	}()
}