package main

import (
	"fmt"
)

var user = ""

func main() {
	throwsPanic(inia)
	test()
}

func throwsPanic(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println(x)
			fmt.Println("recover")
			b = true
		}
	}()

	f()

	fmt.Println("after the run")
	return
}

func inia() {
	defer func() {
		fmt.Println("defer###\n")
	}()
	if user == "" {
		fmt.Print("@@@before panic\n")

		panic("no value for user\n")

		fmt.Print("!!after panic\n")
	}

}
func test() {
	fmt.Println("quit")
}

/*
@@@before panic
defer###

no value for user

recover
quit
*/
