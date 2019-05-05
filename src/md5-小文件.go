package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func main() {
	testFile := os.Args[1]
	file, err := os.Open(testFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	md5h := md5.New()
	io.Copy(md5h, file)
	fmt.Printf("%x", md5h.Sum([]byte("")))
}
