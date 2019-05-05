package main

import (
	"encoding/json"
	"fmt"
)

type A struct {
	B *string `json:"b,omitempty"`
	C int  `json:"c,omitempty"`
}

func main() {
	bb:="2"
	a := A{
		C: 1,
		B: &bb,
	}
	b, _ := json.Marshal(a)
	fmt.Println(string(b))
	fmt.Println(*a.B)


	var s A 
	str:=`{"c":1}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s)
	//fmt.Println(*s.B)
}


//如果不用
// B *string `json:"b,omitempty"`
// 这种方式
// 就会造成有string的默认值

// 如果是
// type A struct {
// 	B string `json:"b,omitempty"`
// 	C int  `json:"c,omitempty"`
// }

// 	a := A{
// 		C: 1,
// 		B: "",
// 	}

// =>
// {"c":1}
