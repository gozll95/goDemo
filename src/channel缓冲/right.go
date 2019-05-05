package main

import "fmt"

func Afunction(ch chan int) {
	fmt.Println("finish")
	<-ch
}

func main() {
	ch := make(chan int) //无缓冲的channel
	go Afunction(ch)
	ch <- 1
	//输出结果
	//finish
}

/*
第一段代码：

        1. 创建了一个无缓冲channel

        2. 启动了一个goroutine，这个routine中对channel执行取出操作，但是因为这时候channel为空，所以这个取出操作发生阻塞，但是主routine可没有发生阻塞，它还在继续运行呢

        3. 主goroutine这时候继续执行下一行，往channel中放入了一个数据

        4. 这时阻塞的那个routine检测到了channel中存在数据了，所以接触阻塞，从channel中取出数据，程序就此完毕


*/
