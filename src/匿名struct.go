package main

import (
	"fmt"
)

func main() {
	a := struct {
		Name string
		Age  int
	}{
		Name: "joe",
		Age:  19,
	}
	fmt.Println(a)
}
