/*
使用sync.Once对象可以使得函数多次调用只执行一次，其结构为:
type Once struct{
	m Mutex
	done unit32
}

func(o *Once)Do(f func())

用done来记录执行次数,用m来保证仅被执行一次。只有一个Do方法,调用执行。
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

//result:
//Only once
