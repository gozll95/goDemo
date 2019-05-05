package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":7777")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	// read from the connection
	var buf = make([]byte, 10)
	log.Println("start to read from conn")
	n, err := c.Read(buf)
	if err != nil {
		log.Println("conn read error:", err)
	} else {
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}

	n, err = c.Write(buf)
	if err != nil {
		log.Println("conn write error:", err)
	} else {
		log.Printf("write %d bytes, content is %s\n", n, string(buf[:n]))
	}
}
