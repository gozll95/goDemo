package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

func main() {
	//定义俩个管道，任务channel和结果channel
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	//启动3个协程，因为jobs里面没有数据，他们都会阻塞。
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	//发送9个任务到jobs channel，然后关闭管道channel。
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	//然后收集所有的结果。
	for a := 1; a <= 9; a++ {
		<-results
	}
}

/*
worker 3 processing job 1
worker 1 processing job 2
worker 2 processing job 3
worker 2 processing job 4
worker 1 processing job 5
worker 3 processing job 6
worker 2 processing job 7
worker 3 processing job 9
worker 1 processing job 8
*/
