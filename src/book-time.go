package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now())
	a := time.NewTimer(4 * time.Second)
	time.Sleep(5 * time.Second)
	fmt.Println(time.Now())
	a.Reset(4 * time.Second)
	//b := <-a.C
//	fmt.Println("time.C:", b)
	time.Sleep(5 * time.Second)
	d := <-a.C
	fmt.Println("time.C:", d)
}

/*
2018-01-10 13:11:49.957261727 +0800 CST
2018-01-10 13:11:54.962358924 +0800 CST
time.C: 2018-01-10 13:11:53.959163413 +0800 CST
time.C: 2018-01-10 13:11:58.96308653 +0800 CST

如果定时器到期了,由于time.C是缓冲为1的，如果没有及时拿数据，会一直存着数据。
*/

