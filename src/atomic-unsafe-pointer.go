package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Expensive struct {
	Data int
}

var wg sync.WaitGroup

var instance unsafe.Pointer

var mutex = &sync.Mutex{}

func GetInstance() *Expensive {
	defer wg.Done()
	if inst := (*Expensive)(atomic.LoadPointer(&instance)); inst != nil {
		return inst
	} else {
		//mutex.Lock()
		//defer mutex.Unlock()

		//time.Sleep(5 * time.Second)
		inst = &Expensive{42}
		atomic.StorePointer(&instance, (unsafe.Pointer)(inst))

		return (*Expensive)(instance)
	}

}

func Broken() {
	wg.Add(1)
	go GetInstance()
	wg.Add(1)
	go GetInstance() // may also take 5+ seconds
}

func main() {
	Broken()
	wg.Wait()
	fmt.Println(instance)
}
