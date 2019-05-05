package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.14

	v := reflect.ValueOf(x)

	fmt.Println("settability of v:", v.CanSet())
}

// settability of v: false

// 这段代码中传给reflect.ValueOf的是x的一个副本，而不是x本身，那么v就是不可以被Set的，
// 所以通过SetFloat方法是不被允许的。那么如何才能被允许被Set呢，也许有人会想到了对于函数，
// 我们可以通过传给函数一个指向参数的指针来达到修改参数本身的作用，同样的道理，这里也可以通过传指针：
