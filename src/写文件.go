package main

import (
	"fmt"
	"os"
)

func main() {
	userFile := "astaxie.txt"
	fout, err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	teststring := "xxxxxxxxxx"
	defer fout.Close()
	for i := 0; i < 10; i++ {
		fout.WriteString("Just a test" + teststring + "\r\n")
		fout.Write([]byte("Just a test\r\n"))
	}
}
