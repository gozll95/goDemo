package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
	Age  int
}

func main() {
	u := User{"张三", 20}

	t := reflect.TypeOf(u)
	fmt.Println(t)

	v := reflect.ValueOf(u)
	fmt.Println(v)

	u1 := v.Interface().(User)
	fmt.Println(u1)

	t1 := v.Type()
	fmt.Println(t1)

	fmt.Println(t.Kind())

	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Name)
	}
	for i := 0; i < t.NumMethod(); i++ {
		fmt.Println(t.Method(i).Name)
	}

	mPrint := v.MethodByName("Print")
	args := []reflect.Value{reflect.ValueOf("前缀")}
	fmt.Println(mPrint.Call(args))
}
func (u User) Print(prfix string) {
	fmt.Printf("%s:Name is %s,Age is %d", prfix, u.Name, u.Age)
}

/*
main.User
{张三 20}
{张三 20}
main.User
struct
Name
Age
Print
前缀:Name is 张三,Age is 20[]
*/
