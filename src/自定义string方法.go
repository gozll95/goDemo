package main

import (
	"fmt"
	"strconv"
)

type Person struct {
	name string
	age  int
}

func (p Person) String() string {
	return "(name: " + p.name + " - age: " + strconv.Itoa(p.age) + " years)"
}

func main() {
	p := Person{"Dennis", 70}
	fmt.Println("Person is", p)
}
