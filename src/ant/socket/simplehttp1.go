package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]

	conn, err := net.Dial("tcp", service)
	checkError(err)

	_, err = conn.Write([]byte("HEAD /HTTP/1.0\r\n\r\n"))
	checkError(err)

	result, err := readFully(conn)
	checkError(err)

	fmt.Println(string(result))

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result.Bytes(), nil
}

/*
 $ go build simplehttp.go
    $ ./simplehttp zhu.me:80
    HTTP/1.1 301 Moved Permanently
    Server: nginx/1.0.14
    Date: Mon, 21 May 2012 03:15:08 GMT
    Content-Type: text/html
    Content-Length: 184
    Connection: close
    Location: https://zhu.me/
*/

/*
实际上，Dial()函数是对DialTCP()、DialUDP()、DialIP()和DialUnix()的  。我
们也可以直接调用这些函数，它们的功能是一 的。这些函数的原型如下:
func DialTCP(net string, laddr, raddr *TCPAddr) (c *TCPConn, err error) func DialUDP(net string, laddr, raddr *UDPAddr) (c *UDPConn, err error) func DialIP(netProto string, laddr, raddr *IPAddr) (*IPConn, error)
func DialUnix(net string, laddr, raddr *UnixAddr) (c *UnixConn, err error)
*/
