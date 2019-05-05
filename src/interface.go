package main

import (
	"fmt"
	"net"
)

func main() {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			fmt.Println(addr)
			fmt.Println(*net.IPNet)
			if _, ok := addr.(*net.IPNet); ok {
				fmt.Println("heh")
			}
		}
	}
}
