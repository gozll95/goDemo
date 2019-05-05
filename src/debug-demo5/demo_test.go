package main

import (
	"fmt"
	"sync"
	"testing"
)

func Test_Demo(t *testing.T) {
	go teller()
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleVistor()
		}()
	}
	vistor := getVistors()
	fmt.Println(vistor)
}
