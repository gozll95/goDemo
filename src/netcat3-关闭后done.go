//Netcat is a simple read/write client for TCP servers.

package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	con := conn.(*net.TCPConn)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		log.Println("start")
		io.Copy(os.Stdout, con) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(con, os.Stdin)
	//con.CloseWrite()
	con.CloseRead()
	//conn.Close()
	<-done // wait for background goroutine to finish
	log.Println("quit")
}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

/*
当用户关闭了标准输入，主goroutine中的mustCopy函数调用将返回，然后调用conn.Close()关闭读和写方向的网络连接。关闭网络链接中的写方向的链接将导致server程序收到一个文件（end-of-ﬁle）结束的信号。关闭网络链接中读方向的链接将导致后台goroutine的io.Copy函数调用返回一个“read from closed connection”（“从关闭的链接读”）类似的错误，因此我们临时移除了错误日志语句；在练习8.3将会提供一个更好的解决方案。（需要注意的是go语句调用了一个函数字面量，这Go语言中启动goroutine常用的形式。）

在后台goroutine返回之前，它先打印一个日志信息，然后向done对应的channel发送一个值。主goroutine在退出前先等待从done对应的channel接收一个值。因此，总是可以在程序退出前正确输出“done”消息。
*/
