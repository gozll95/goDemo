/*
 Echo 服务

编写一个简单的 echo 服务。使其监听于本地的 TCP 端口 8053 上。它应当可以读取一行（以换行符结尾），将这行原样返回然后关闭连接。
让这个服务可以并发，这样每个请求都可以在独立的 goroutine 中进行处理。

*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)

	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	_, err = conn.Write([]byte("-------\n" + line + " u are right\n"))

	if err != nil {
		return
	}

}
