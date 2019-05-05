package main

import (
	"log"
	"net"
	"time"
)

func establishConn(i int) net.Conn {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Printf("%d: dial error: %s", i, err)
		return nil
	}
	log.Println(i, ":connect to server ok")
	return conn
}

func main() {
	var sl []net.Conn
	for i := 1; i < 1000; i++ {
		conn := establishConn(i)
		if conn != nil {
			sl = append(sl, conn)
		}
	}

	time.Sleep(time.Second * 10000)
}

/*
2018/05/11 17:42:30 99 :connect to server ok
2018/05/11 17:42:30 100 :connect to server ok
2018/05/11 17:42:30 101 :connect to server ok
2018/05/11 17:42:30 102 :connect to server ok
2018/05/11 17:42:30 103 :connect to server ok
2018/05/11 17:42:30 104 :connect to server ok
2018/05/11 17:42:30 105 :connect to server ok
2018/05/11 17:42:30 106 :connect to server ok
2018/05/11 17:42:30 107 :connect to server ok
2018/05/11 17:42:30 108 :connect to server ok
2018/05/11 17:42:30 109 :connect to server ok
2018/05/11 17:42:30 110 :connect to server ok
2018/05/11 17:42:30 111 :connect to server ok
2018/05/11 17:42:30 112 :connect to server ok
2018/05/11 17:42:30 113 :connect to server ok
2018/05/11 17:42:30 114 :connect to server ok
2018/05/11 17:42:30 115 :connect to server ok
2018/05/11 17:42:30 116 :connect to server ok
2018/05/11 17:42:30 117 :connect to server ok
2018/05/11 17:42:30 118 :connect to server ok
2018/05/11 17:42:30 119 :connect to server ok
2018/05/11 17:42:30 120 :connect to server ok
2018/05/11 17:42:30 121 :connect to server ok
2018/05/11 17:42:30 122 :connect to server ok
2018/05/11 17:42:30 123 :connect to server ok
2018/05/11 17:42:30 124 :connect to server ok
2018/05/11 17:42:30 125 :connect to server ok
2018/05/11 17:42:30 126 :connect to server ok
2018/05/11 17:42:30 127 :connect to server ok
2018/05/11 17:42:30 128 :connect to server ok
2018/05/11 17:42:30 129 :connect to server ok // 上面会一次性打满server的backlog
2018/05/11 17:42:43 130 :connect to server ok // 然后会开始 每10s被accept
2018/05/11 17:42:57 131 :connect to server ok
*/
