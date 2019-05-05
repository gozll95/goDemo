package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":7777")
	if err != nil {
		log.Println("dial error:", err)
		return
	}

	conn.Close()
	log.Println("close ok")

	var buf = make([]byte, 32)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("read error:", err)
	} else {
		log.Printf("read % bytes, content is %s\n", n, string(buf[:n]))
	}

	n, err = conn.Write(buf)
	if err != nil {
		log.Println("write error:", err)
	} else {
		log.Printf("write % bytes, content is %s\n", n, string(buf[:n]))
	}

	time.Sleep(time.Second * 1000)
}

/*
flower@:~/workspace/learngo/src/myGoNotes/TCP情况分析/关闭连接$ go run server.go
2017/08/04 21:12:19 start to read from conn
2017/08/04 21:12:19 conn read error: EOF
2017/08/04 21:12:19 write 10 bytes, content is


flower@:~/workspace/learngo/src/myGoNotes/TCP情况分析/关闭连接$ go run client.go
2017/08/04 21:12:19 begin dial...
2017/08/04 21:12:19 close ok
2017/08/04 21:12:19 read error: read tcp 127.0.0.1:55933->127.0.0.1:7777: use of closed network connection
2017/08/04 21:12:19 write error: write tcp 127.0.0.1:55933->127.0.0.1:7777: use of closed network connection
*/

/*
从client1的结果来看，在己方已经关闭的socket上再进行read和write操作，会得到”use of closed network connection” error；
从server1的执行结果来看，在对方关闭的socket上执行read操作会得到EOF error，但write操作会成功，因为数据会成功写入己方的内核socket缓冲区中，即便最终发不到对方socket缓冲区了，因为己方socket并未关闭。因此当发现对方socket关闭后，己方应该正确合理处理自己的socket，再继续write已经无任何意义了。
*/
