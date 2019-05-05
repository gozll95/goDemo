package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

	visit := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Println("dir:  ", path)
		} else {
			fmt.Println("file: ", path)
		}
		return nil
	}

	err := filepath.Walk(flag.Arg(0), visit)
	if err != nil {
		log.Fatal(err)
	}
}
