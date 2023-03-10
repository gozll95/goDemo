package main

import (
	"sync"
)

type Counter interface {
	Add(int) int
}

type countService struct {
	v  int
	mu sync.Mutex
}

func (c *countService) Add(v int) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v += v
	return c.v
}
