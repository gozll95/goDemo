package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	userFile := "astaxie.txt"
	dat, err := ioutil.ReadFile(userFile)
	check(err)
	fmt.Println(string(dat))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
