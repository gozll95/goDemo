package main

import (
	"fmt"
	"reflect"
)

type T struct {
	A int
	B string
}

func main() {
	t := T{23, "hello world"}
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

// 0: A int = 23
// 1: B string = hello world

// 同样可以通过以下类似的代码来修改T的值：

// s.Field(0).SetInt(22)
// s.Field(1).SetString("XXOO")

// 好了，反射就基本如此吧，记住三点即可：

// 1. Reflection goes from interface value to reflection Object.

// 2. Reflection goes from refelction object to interface value.

// 3. To modify a reflection object, the value must be settable.
