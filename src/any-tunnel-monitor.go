package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	// 域名解析
	ns, err := net.LookupHost("www.baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(ns)

	// 测试端口连通性
	ip := net.ParseIP(ns[0])
	port := 80
	tcpaddr := &net.TCPAddr{IP: ip, Port: port}
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v:%v connection ok", ip, port)
	conn.Close()

	// 测试业务连通性
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://xxxx/download", nil)
	if err != nil {
		panic(err)
	}

	

}
