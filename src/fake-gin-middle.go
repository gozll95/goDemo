package main

import (
	"fmt"
	"math"
)

const (
	defaultMemory      = 32 << 20 // 32 MB
	abortIndex    int8 = math.MaxInt8 / 2
)

type errorMsgs []string

type Context struct {
	index    int8
	handlers []HandlerFunc
	Errors   errorMsgs
}

func (c *Context) reset() {
	c.handlers = nil
	c.index = -1
	c.Errors = c.Errors[0:0]
}

func (c *Context) Error(err error) string {
	c.Errors = append(c.Errors, err.Error())
	return "has error"
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func main() {
	c := &Context{}
	c.reset()
	c.handlers = []HandlerFunc{A(), B(), C()}
	c.Next()
}

type HandlerFunc func(c *Context)

func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func A() func(c *Context) {
	return func(c *Context) {
		fmt.Println("a")

		if 1==1{
			c.Abort()
			return 
		}
		fmt.Println("aa")

		c.Next()
	}
}

func B() func(c *Context) {
	return func(c *Context) {
		fmt.Println("b")
		c.Next()
	}
}

func C() func(c *Context) {
	return func(c *Context) {
		fmt.Println("c")
		c.Next()
	}
}
