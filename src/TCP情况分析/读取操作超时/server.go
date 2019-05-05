package main

import (
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Println("errorlisten:", err)
		return
	}
	defer l.Close()
	log.Println("listen ok")

	var i int
	for {
		time.Sleep(time.Second * 10)
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}
		i++
		log.Printf("%d: accept a new connection\n", i)
		go handleConn(c)

	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		time.Sleep(10 * time.Second)
		var buf = make([]byte, 65536)
		log.Println("start to read from conn")
		c.SetReadDeadline(time.Now().Add(time.Microsecond * 10))
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}
