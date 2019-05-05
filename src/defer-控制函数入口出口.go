package main

import (
	"fmt"
	"log"
	"time"
)

func bigSlowOperation() {
	defer trace("bigSlowOperation")()
	fmt.Println("start")
	time.Sleep(5 * time.Second)
	fmt.Println("end")
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() {
		log.Printf("exit %s (%s)", msg, time.Since(start))
	}
}

func main() {
	bigSlowOperation()
}

/*
bigSlowOperation被调时，trace会返回一个函数值，该函数值会在bigSlowOperation退出时被调用。通过这种方式， 我们可以只通过一条语句控制函数的入口和所有的出口，甚至可以记录函数的运行时间，如例子中的start。需要注意一点：不要忘记defer语句后的圆括号，否则本该在进入时执行的操作会在退出时执行，而本该在退出时执行的，永远不会被执行。
*/

/*
flower@:~/workspace/learngo/src/myGoNotes$ go run defer-控制函数入口出口.go
2017/11/12 16:41:07 enter bigSlowOperation
start
end
2017/11/12 16:41:17 exit bigSlowOperation (10.0013568s)
*/
