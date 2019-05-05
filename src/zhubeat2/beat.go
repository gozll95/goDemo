package main

import (
	"fmt"
	"strconv"
	"time"
	"zhubeat/beat"
	"zhubeat/beat/output"
	"zhubeat/beat/queue"
	"zhubeat/beat/worker"
	"zhubeat/lib/manager"
)

func main() {
	maxWorkers := 2
	maxQueue := 10

	var jobChan chan manager.Job
	jobChan = make(chan manager.Job, maxQueue)

	// dispather
	dispatcher := manager.NewDispatcher(maxWorkers, jobChan, worker.NewWorker)
	dispatcher.Init()
	dispatcher.Start()

	// output
	hosts := []string{"127.0.0.1:3111", "127.0.0.1:7700", "127.0.0.1:6000"}
	timeout := 10 * time.Second
	ttl := 500 * time.Millisecond
	concurrency := 10
	outputManager, err := output.NewOutPutManager("tcp", hosts, timeout, ttl, concurrency, queue.ToOutputQueue)
	if err != nil {
		panic(err)
	}
	outputManager.Run()

	// job里放batch
	go func() {
		for i := 0; i < 100000; i++ {
			jobChan <- manager.Job(beat.Job(strconv.Itoa(i) + "\n"))
		}
	}()

	time.Sleep(100 * time.Second)

	dispatcher.Close()
	fmt.Println("xxx")
	outputManager.Stop()
	fmt.Println("yyy")

	fmt.Println(queue.ToOutputQueue.Dump())

	//test()

}

// func test() {

// 	hosts := []string{"127.0.0.1:8000", "127.0.0.1:7000", "127.0.0.1:6000"}
// 	conn, err := transport.NewClient("tcp", hosts, time.Duration(100*time.Millisecond))
// 	if err != nil {
// 		panic(err)
// 	}

// 	go func() {
// 		var jobs []manager.Job
// 		for i := 0; i < 100; i++ {
// 			jobs = append(jobs, beat.Job("from client1:"+strconv.Itoa(i)+"\n"))
// 		}

// 		outClient := beat.NewOutputClient(conn, 500*time.Millisecond)
// 		err = outClient.Publish(jobs)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	go func() {
// 		var jobs []manager.Job
// 		for i := 0; i < 100; i++ {
// 			jobs = append(jobs, beat.Job("from client2:"+strconv.Itoa(i)+"\n"))
// 		}

// 		outClient := beat.NewOutputClient(conn, 500*time.Millisecond)

// 		err = outClient.Publish(jobs)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	time.Sleep(100 * time.Second)
// 	// for {
// 	// 	data := time.Now().Format("15:04:05\n")
// 	// 	_, err = outClient.Client.Write([]byte(data))
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	time.Sleep(1 * time.Second)
// 	// }

// }

func test2() {

}
