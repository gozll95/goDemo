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
		time.Sleep(5 * time.Second)
		var buf = make([]byte, 60000)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
		}

		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

/*
在Darwin上，这个size大约在679468bytes。后续当server5每隔5s进行Read时，OS socket缓冲区腾出了空间，client5就又可以写入了：
*/
