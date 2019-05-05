package main

import (
	"fmt"
	"reflect"
)

type MyInt int

func main() {
	var x MyInt = 7
	v := reflect.ValueOf(x)
	fmt.Println("type:", v.Type())
	fmt.Println("kind is uint8:", v.Kind())

}

// 2.如果Kind方法是描述相关的类型，而不是静态的类型，例如用户自定义了一个类型：

// type: main.MyInt
// kind is uint8: in
