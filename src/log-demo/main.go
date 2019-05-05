package main

import (
	"log"
	"log-demo/test1"
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	test1.Test()
	test1.TestLog()
	test1.TestLog4()
}
