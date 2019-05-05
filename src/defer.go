package main

import (
	"fmt"
)

func main() {
	a, b, c := 1, 2, 3
	fmt.Println(a)
	defer fmt.Println(b)
	defer fmt.Println(c)

	for i := 0; i < 3; i++ {
		fmt.Println("hehe")
		defer fmt.Println(i)
		fmt.Println("xx")
	}
}

//defer是先return后defer

/*
a
hehe
xx
hehe
xx
hehe
xx
2
1
0
c
b
*/
/*
fmt.Println("hehe")
defer fmt.Println(i)-->0
fmt.Println("xx")
fmt.Println("hehe")
defer fmt.Println(i)-->1
fmt.Println("xx")
fmt.Println("hehe")
defer fmt.Println(i)-->2
fmt.Println("xx")
-->
fmt.Println("hehe")
fmt.Println("xx")
fmt.Println("hehe")
fmt.Println("xx")
fmt.Println("hehe")
fmt.Println("xx")
defer fmt.Println(i)-->2
defer fmt.Println(i)-->1
defer fmt.Println(i)-->0
*/
