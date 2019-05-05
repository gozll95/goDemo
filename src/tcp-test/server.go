package main

import (
	"log"
	"net"
	"time"
)

//!+
func handleConn(c net.Conn) {
	defer c.Close()
	time.Sleep(100 * time.Second)
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
