package main

import (
	"log"
	"net"
)

func main() {

}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		var buf = make([]byte, 10)
		n, err := c.Read(buf)
		if err != nil {
			return
		}
		log.Printf("read %d bytes, content is %s\n", s, string(buf[:n]))
	}

}
