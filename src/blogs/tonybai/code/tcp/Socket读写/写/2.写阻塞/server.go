//server.go

package main

import (
	"log"
	"net"
	"time"
)

func handleConn(c net.Conn) {
	defer c.Close()
	time.Sleep(time.Second * 10)
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
			break
		}

		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
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
TCP连接通信两端的OS都会为该连接保留数据缓冲,一端调用Write后,实际上数据是写入到OS的协议栈的数据缓冲的。
TCP是全双工通信,因此每个方向都有独立的数据缓冲。当发送方将对方的接收缓冲区以及自身的发送缓冲区写满后,Write
就会阻塞。
*/

/*
Server5在前10s中并不Read数据，因此当client5一直尝试写入时，写到一定量后就会发生阻塞：

$go run client5.go

2015/11/17 14:57:33 begin dial...
2015/11/17 14:57:33 dial ok
2015/11/17 14:57:33 write 65536 bytes this time, 65536 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 131072 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 196608 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 262144 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 327680 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 393216 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 458752 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 524288 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 589824 bytes in total
2015/11/17 14:57:33 write 65536 bytes this time, 655360 bytes in total


在Darwin上，这个size大约在679468bytes。后续当server5每隔5s进行Read时，OS socket缓冲区腾出了空间，client5就又可以写入了：
$go run server5.go
2015/11/17 15:07:01 accept a new connection
2015/11/17 15:07:16 start to read from conn
2015/11/17 15:07:16 read 60000 bytes, content is
2015/11/17 15:07:21 start to read from conn
2015/11/17 15:07:21 read 60000 bytes, content is
2015/11/17 15:07:26 start to read from conn
2015/11/17 15:07:26 read 60000 bytes, content is
....

client端：

2015/11/17 15:07:01 write 65536 bytes this time, 720896 bytes in total
2015/11/17 15:07:06 write 65536 bytes this time, 786432 bytes in total
2015/11/17 15:07:16 write 65536 bytes this time, 851968 bytes in total
2015/11/17 15:07:16 write 65536 bytes this time, 917504 bytes in total
2015/11/17 15:07:27 write 65536 bytes this time, 983040 bytes in total
2015/11/17 15:07:27 write 65536 bytes this time, 1048576 bytes in total
.... ...
*/
