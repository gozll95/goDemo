package client

import (
	"errors"
	"fmt"
	"io"
	"net"
	"npool/gopool"
	"time"
)

var addrs []string = []string{"127.0.0.1:8000", "127.0.0.1:7000", "127.0.0.1:6000", "127.0.0.1:5000"}

type Client struct {
	name        string
	pool        *gopool.ChanConnPool
	conntimeout time.Duration
	retry       int
	conn        net.Conn
}

func NewClient(name string) (*Client, error) {
	epool, err := gopool.NewChanConnPool(&gopool.ConnPoolReq{
		Addrs: addrs,
		Create: func(addr string, timeout time.Duration) (interface{}, error) {
			cli, err := net.DialTimeout("tcp", addr, timeout)
			return cli, err
		},
		IsOpen: func(cli interface{}) bool {
			if cli != nil {
				return true
			}
			return false
		},
		Down: func(cli interface{}) {
			c := cli.(net.Conn)
			c.Close()
		},
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		name:        name,
		pool:        epool,
		conntimeout: 2 * time.Second,
		retry:       2,
	}, nil
}

func (c *Client) get() (conn net.Conn, err error) {
	var ok bool
	cli, err := c.pool.Get()
	if err != nil {
		return
	}
	if conn, ok = cli.(net.Conn); !ok {
		err = errors.New("cli is nil.")
	}
	return
}

func (c *Client) put(conn net.Conn, err *error) {
	safe := false
	if *err == nil {
		safe = true
	}
	c.pool.Put(conn, safe)
}

// reget
func (c *Client) reget(src net.Conn, srcErr *error) (conn net.Conn, err error) {
	c.put(src, srcErr)
	return c.get()
}

func (c *Client) GetConnCount() map[string]int {
	return c.pool.GetConnCount()
}

func (c *Client) GetHealthy() map[string]bool {
	return c.pool.GetHealthy()
}

func (c *Client) Write(messages []string) (err error) {
	var (
		conn       net.Conn
		needReconn bool
		firstConn  bool
	)

	defer func() {
		c.put(conn, &err)
	}()

	firstConn = true

	for _, message := range messages {
		time.Sleep(1 * time.Second)
		if firstConn {
			conn, err = c.get()
			if err != nil {
				fmt.Println("store message: ", message)
				needReconn = true
				continue
			}
			firstConn = false
		}
		if needReconn {
			conn, err = c.reget(conn, &err)
			if err != nil {
				fmt.Println("store message: ", message)
				continue
			}
		}
		_, err = io.WriteString(conn, "from:"+c.name+"message:"+message+"\n")
		if err != nil {
			// store
			fmt.Println("store message: ", message)
			needReconn = true
		} else {
			needReconn = false
		}

	}

	return err
}
