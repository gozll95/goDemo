package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.14

	p := reflect.ValueOf(&x)

	fmt.Println("type of p:", p.Type())
	fmt.Println("settability of p:", p.CanSet())
}

// type of p: *float64
// settability of p: false

// 我们可以看到p的类型是*float64，而不是float64了，但是为什么还是不可以被Set呢,因为这里p是一个指针，我们并不是要Set这个指针的值，
// 而是要Set指针所指内容的值（也就是*p）,所以这里p仍然是不可被Set的，
// 我们可以通过Value的Elem方法来指针所指向内容的Value：
