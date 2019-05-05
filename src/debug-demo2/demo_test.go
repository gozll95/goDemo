package main

import (
	"sync"
	"testing"
)

func Test_Demo(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		go handleVistors()
	}
}
