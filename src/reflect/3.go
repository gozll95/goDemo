// 1.大的数据类型可以包含小的数据类型，例如int64可以是任意的整数(int8,uint8,int32等),但是需要转换一下。

package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x uint8 = 'x'
	v := reflect.ValueOf(x)
	fmt.Println("type:", v.Type())
	fmt.Println("kind is uint8:", v.Kind() == reflect.Uint8)
	x = uint8(v.Uint())
}

// type: uint8
// kind is uint8: true
