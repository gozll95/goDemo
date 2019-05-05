package main

import (
	"fmt"
	"time"
)

var mapChan = make(chan map[string]int, 1)

func main() {
	syncChan := make(chan struct{}, 2)
	go func() {
		for {
			if elem, ok := <-mapChan; ok {
				elem["count"]++
			} else {
				break
			}
		}
		fmt.Println("Stopped.[receiver]")
		syncChan <- struct{}{}
	}()
	go func() {
		countMap := make(map[string]int)
		for i := 0; i < 5; i++ {
			mapChan <- countMap
			time.Sleep(time.Microsecond)
			fmt.Printf("The count map:%v.[sender]\n", countMap)
		}
		close(mapChan)
		syncChan <- struct{}{}
	}()
	<-syncChan
	<-syncChan

}

/*
The count map:map[count:1].[sender]
The count map:map[count:2].[sender]
The count map:map[count:3].[sender]
The count map:map[count:4].[sender]
The count map:map[count:5].[sender]
Stopped.[receiver]
*/

//
/*
当接收方从通道接受到一个值类型的值时，对该值的修改就不会影响到发送方持有的那个源值。但是对于引用类型的值来说，这种修改就会同时影响到双方持有的值。
*/