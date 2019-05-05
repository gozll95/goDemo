package main

import "fmt"

type A struct {
	Name string
}

func main() {
	var a = make([]*A, 5)
	var b = make([]int, 5)
	var c = make([]int, 5)
	var d = make([]A, 5)

	for i, _ := range a {
		a[i] = &A{}
		test(a[i])
		// a[i].Name = "aa"
		// fmt.Println(a[i])
	}

	for i, v := range b {
		if i == 0 {
			// b[1] = 1
			// b[2] = 2
			v = 11
		}
		c[i] = v
	}

	for i, _ := range d {
		d[i] = A{}
		//test(d[i])
		d[i].Name = "aa"
	}

	fmt.Println(a[0].Name)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
}

func test(a interface{}) {
	switch v := a.(type) {
	case *A:
		v.Name = "aa"
	case A:
		v.Name = "aa"
		fmt.Println("v is ", v)
	default:
		fmt.Println("haha")
	}
	fmt.Println("a is ", a)
}

func test1(a *A) {
	a.Name = "aa"
}
