package client

import (
	"book-rpc/server"
	"fmt"
	"log"
	"net/rpc"
)

func test() {
	//与RPC服务端建立连接
	client, err := rpc.DialHTTP("tpc", serverAddress+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	//同步调用程序顺序执行
	args := &server.Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

	//异步调用方式
	quotient := new(server.Quotient)
	divCall := client.Go("Arith.Divide", args, &quotient, nil)
	replyCall := <-divCall.Done

}

/*
 * rpc 调用 有同步/异步 两种方式
 */
