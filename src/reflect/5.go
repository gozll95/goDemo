package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.14
	v := reflect.ValueOf(x)

	fmt.Println(v.Interface())
	fmt.Printf("value is %7.1e\n", v.Interface())

	y := v.Interface().(float64)
	fmt.Println(y)

}

// 3.14
// value is 3.1e+00
// 3.14

// Relection goe s from reflection object to interface value
