// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 222.

// Clock is a TCP server that periodically writes the time.
package main

import (
	"bytes"
	"io"
	"log"
	"net"
)

const DELIMITER = '\n'

// read 会从连接中读数据直到遇到参数delim代表的字节。
func read(conn net.Conn) (string, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELIMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		data, err := read(c)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(data)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:6000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen on 6000")
	//!+
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
	//!-
}

/*
for conn=accept go handle conn
handle 里 外层 defer close 里层 for wirte
*/
