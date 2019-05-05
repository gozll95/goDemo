package main

import (
	"fmt"
)

func close(x int) (f func(y int) int) {
	return func(y int) int {
		return x + y
	}
}

func main() {
	f := close(10)
	f1 := f(1)
	f2 := f(2)
	fmt.Println(f1, f2)
}

//func return func 就是 闭包
//闭包先变再变
