package main

import "fmt"

type A struct {
	vex []int
}

func main() {
	a := A{}
	a.vex = []int{}
	for i := 0; i < 7; i++ {
		a.vex[i] = i
	}
	fmt.Println(a)

}
