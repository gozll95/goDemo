package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func establishConn(i int) net.Conn {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		log.Printf("%d: dial error: %s", i, err)
		return nil
	}
	log.Println(i, ":connect to server ok")
	return conn
}

func main() {
	conn := establishConn(0)
	read(conn)
}

func read(conn net.Conn) {
	n, err := ioutil.ReadAll(conn)
	if err != nil {
		return
	}
	fmt.Println(string(n))

}
