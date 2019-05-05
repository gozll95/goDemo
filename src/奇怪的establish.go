package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	URL = "http://localhost:9090/"
)

func PrintLocalDial(network, addr string) (net.Conn, error) {
	fmt.Printf("network is %s\taddr is %s\n", network, addr)
	//network is tcp	addr is localhost:9090
	dial := net.Dialer{
	// Timeout:   30 * time.Second,
	// KeepAlive: 30 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}
	fmt.Println("connect done, use", conn.LocalAddr().String())
	return conn, err
}

func doGet(url string, id int) error {
	transport := http.Transport{
		Dial:              PrintLocalDial,
		DisableKeepAlives: true,
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%d: %s -- %v\n", id, string(buf), err)
	if err := resp.Body.Close(); err != nil {
		fmt.Println(err)
	}
	return err
}
func main() {
	for {
		go doGet(URL, 1)
		go doGet(URL, 2)
		time.Sleep(2 * time.Second)
	}

}
