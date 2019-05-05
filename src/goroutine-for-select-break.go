// http://www.tuicool.com/articles/za6Zbey

/*
简单来说， break 是终止了 switch 的运算，
对于 for 循环来说没有作用，这点对于 for select 也一样，
如果要终止 for 循环，则需要使用label来标示跳出位置。
*/

//A "break" statement terminates execution of the innermost "for", "switch" or "select" statement.

package main

import (
	"fmt"
	"time"
)

var c chan int

func ready(w string, sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
	fmt.Println(w, "is ready!")
	c <- 1
}

func main() {
	c = make(chan int, 2)
	go ready("Tea", 2)
	go ready("Coffee", 1)
	fmt.Println("I'm waiting")

	for i := 0; i < 2; i++ {
		<-c
	}
	/*
		i := 0
	L:
		for {
			select {
			case <-c:
				i++
				fmt.Println(i)
				break L

			}
		}
		// time.Sleep(5 * time.Second)
	*/
}
