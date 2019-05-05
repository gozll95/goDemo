package main

import (
	"fmt"
	"runtime"
)

func main() {
	funcName, file, line, ok := runtime.Caller(0)
	if ok {
		fmt.Println("Func Name=" + runtime.FuncForPC(funcName).Name())
		fmt.Printf("file: %s    line=%d\n", file, line)
	}
}

/*
flower@:~/workspace/learngo/src/myGoNotes$ go run runtime1.go
Func Name=main.main
file: /Users/flower/workspace/learngo/src/myGoNotes/runtime1.go    line=9
*/

/*
哪里调用了runtime就在哪里打印堆栈
比如log是在output方法里调用runtime
*/
