package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.14
	
	fmt.Println("type:", reflect.TypeOf(x))

	fmt.Println("vlaue:", reflect.ValueOf(x))
}

// type: float64
// vlaue: 3.14





1.大的数据类型可以包含小的数据类型，例如int64可以是任意的整数(int8,uint8,int32等),但是需要转换一下。

var x uint8 = 'x'
v := reflect.ValueOf(x)
fmt.Println("type:"，v.Type())
fmt.Println("kind is uint8:",v.Kind() == reflect.Uint8)
x = uint8(v.Uint())
2.如果Kind方法是描述相关的类型，而不是静态的类型，例如用户自定义了一个类型：

type MyInt int
var x MyInt = 7
v := reflect.ValueOf(x)
那么v.Kind()的返回值是reflect.Int，尽管x的静态类型是MyInt而不是int。Kind无法描述一个MyInt的int但是Type可以。