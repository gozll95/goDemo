package main

import (
	"fmt"
	"time"
)

type Job struct {
	JobId int
}

func worker(id int) {
	go func() {
		for {
			fmt.Println("Waiting for job...")
			select {
			// Receive from channel
			case j := <-Jobs:
				if j == nil {
					fmt.Println("Close the worker", id)
					return
				}
				fmt.Println("worker", id, "started  job", j.JobId)
				time.Sleep(time.Second)
				fmt.Println("worker", id, "finished job", j.JobId)
				Results <- true
			}
		}
	}()
}

const channelLength = 3

var (
	Jobs    chan *Job
	Results chan bool
)

func main() {
	Jobs = make(chan *Job, channelLength)
	Results = make(chan bool, channelLength)

	// Start worker goroutines
	for i := 0; i < channelLength; i++ {
		worker(i)
	}

	// Send to channel
	time.Sleep(time.Second)
	for j := 0; j < channelLength; j++ {
		Jobs <- &Job{JobId: j}
	}
	close(Jobs)

	for len(Jobs) != 0 || len(Results) != channelLength {
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Complete main")
}

//close一个管道，判断一个管道是否被关闭，还可以用管道里放指针这种方式来做。
// type Job struct {
// 	JobId int
// }

// Jobs = make(chan *Job,channelLength)
// Jobs ==nil 的时候即使被关闭了

/*
http://www.jianshu.com/p/316af7d6730c

golang中启动一个协程不会消耗太多资源，有人认为可以不用协程池。但是当访问量增大时，可能造成内存消耗完，程序崩溃。于是写了一个协程池的Demo。

Demo中有worker和job。worker是一个协程，在worker中完成一个job。Jobs是一个channel，使用Jobs记录job。当生成一个新任务，就发送到Jobs中。程序启动时，首先启动3个worker协程，每个协程都尝试从Jobs中接收job。如果Jobs中没有job，worker协程就等待。

基本逻辑如下：

Jobs管道存放job，Results管道存放结果。
程序一启动，启动3个worker协程，等待从Jobs管道中取数据。
向Jobs管道中发送3个数据。
关闭Jobs管道。
worker协程从Jobs管道中接收到数据以后，执行程序，把结果放到Results管道中。然后继续等待。
当Jobs管道中没有数据，并且Results有3个数据时。退出主程序。

*/
