package main

import "fmt"

type I interface {
	A() int
	B() string
}

func main() {
	mockEncoder := struct {
		I
	}{}
	fmt.Println(mockEncoder.A())
	fmt.Println(mockEncoder.B())

}
