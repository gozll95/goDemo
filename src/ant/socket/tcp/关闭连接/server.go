//server.go

package main

import (
	"log"
	"net"
)

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

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		log.Println("accept a new connection")
		go handleConn(c)
	}
}

/*

和前面的方法相比，关闭连接算是最简单的操作了。由于socket是全双工的，client和server端在己方已关闭的socket和对方关闭的socket上操作的结果有不同。看下面例子：


$go run server1.go
2015/11/17 17:00:51 accept a new connection
2015/11/17 17:00:51 start to read from conn
2015/11/17 17:00:51 conn read error: EOF
2015/11/17 17:00:51 write 10 bytes, content is

$go run client1.go
2015/11/17 17:00:51 begin dial...
2015/11/17 17:00:51 close ok
2015/11/17 17:00:51 read error: read tcp 127.0.0.1:64195->127.0.0.1:8888: use of closed network connection
2015/11/17 17:00:51 write error: write tcp 127.0.0.1:64195->127.0.0.1:8888: use of closed network connection

从client1的结果来看，在己方已经关闭的socket上再进行read和write操作，会得到”use of closed network connection” error；
从server1的执行结果来看，在对方关闭的socket上执行read操作会得到EOF error，但write操作会成功，因为数据会成功写入己方的内核socket缓冲区中，即便最终发不到对方socket缓冲区了，因为己方socket并未关闭。因此当发现对方socket关闭后，己方应该正确合理处理自己的socket，再继续write已经无任何意义了。
*/
