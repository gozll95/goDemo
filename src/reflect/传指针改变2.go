package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.14

	p := reflect.ValueOf(&x)

	v := p.Elem()

	fmt.Println("settability of v:", v.CanSet(), v.Kind())

	v.SetFloat(2.8)

	fmt.Println(v.Interface())

	fmt.Println(x)
}

// settability of v: true
// 2.8
// 2.8
