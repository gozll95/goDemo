package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Id   int
	Name string
	Age  int
}

func set(o interface{}) {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() {
		fmt.Println("xxx")
		return
	} else {
		v = v.Elem()
	}

	if f := v.FieldByName("Name"); f.Kind() == reflect.String {
		f.SetString("BAYBAY")
	}

}
func main() {
	x := 123
	v := reflect.ValueOf(&x)
	v.Elem().SetInt(999)
	fmt.Println(x)

	u := User{1, "ok", 12}
	set(&u)
	fmt.Println(u)
}
