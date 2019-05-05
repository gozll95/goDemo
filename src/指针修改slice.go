package main

import "fmt"

func main() {
	var a []int
	a = []int{1}
	fmt.Println(a)
	test(&a)
	fmt.Println(a)
}

func test(aa *[]int) {
	*aa = append(*aa, 2)
}
