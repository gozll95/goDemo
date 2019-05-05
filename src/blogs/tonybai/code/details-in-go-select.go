package main

import (
	"fmt"
	"time"
)

func takeARecvChannel() chan int {
	fmt.Println("invoke takeARecvChannel")
	c := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		c <- 1
	}()
	return c
}
func getAStorageArr() *[5]int {
	fmt.Println("invoke getAStorageArr")
	var a [5]int
	return &a
}
func takeASendChannel() chan int {
	fmt.Println("invoke takeASendChannel")
	return make(chan int)
}
func getANumToChannel() int {
	fmt.Println("invoke getANumToChannel")
	return 2
}
func main() {
	select {
	//recv channels
	case (getAStorageArr())[0] = <-takeARecvChannel():
		fmt.Println("recv something from a recv channel")
		//send channels
	case takeASendChannel() <- getANumToChannel():
		fmt.Println("send something to a send channel")
	}
}

/*
运行结果：
$go run select.go
invoke takeARecvChannel
invoke takeASendChannel
invoke getANumToChannel
invoke getAStorageArr
recv something from a recv channel

通过例子我们可以看出：
1) select执行开始时，首先所有case expression的表达式都会被求值一遍，按语法先后次序。
invoke takeARecvChannel
invoke takeASendChannel
invoke getANumToChannel
例外的是recv channel的位于赋值等号左边的表达式（这里是：(getAStorageArr())[0]）不会被求值。
2) 如果选择要执行的case是一个recv channel，那么它的赋值等号左边的表达式会被求值：如例子中当goroutine 3s后向recvchan写入一个int值后，select选择了recv channel执行，此时对=左侧的表达式 (getAStorageArr())[0] 开始求值，输出“invoke getAStorageArr”。
*/