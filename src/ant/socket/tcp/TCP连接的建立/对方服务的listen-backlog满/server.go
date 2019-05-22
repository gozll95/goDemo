package main

import (
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("error listen:", err)
		return
	}
	defer l.Close()
	log.Println("listen ok")

	var i int
	for {
		time.Sleep(time.Second * 10)
		if _, err := l.Accept(); err != nil {
			log.Println("accept error:", err)
			break
		}
		i++
		log.Println("%d: accept a new connection\n", i)
	}
}

// 连接会------> backlog ------->accept 每次从backlog里拿一个conn
// backlog相当于缓冲层

/*
2018/05/11 17:42:11 listen ok
2018/05/11 17:42:30 %d: accept a new connection
 1
2018/05/11 17:42:40 %d: accept a new connection
 2
2018/05/11 17:42:50 %d: accept a new connection
 3
2018/05/11 17:43:00 %d: accept a new connection
 4
^Csignal: interrupt
*/
