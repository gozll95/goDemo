package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"video-code/rpc"
)

func main() {
	rpc.Register(rpcdemo.DemoService{})
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}

/*
➜ /Users/flower >telnet localhost 1234
Trying ::1...
Connected to localhost.
Escape character is '^]'.
{"method":"abc.def"}
{"id":null,"result":null,"error":"rpc: can't find service abc.def"}
{"method":"DemoService.Div","params":[{"A":3,"B":2}],"id":1}
{"id":1,"result":1.5,"error":null}
*/
