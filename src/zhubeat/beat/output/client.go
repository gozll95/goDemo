package output

import (
	"errors"
	"fmt"
	"net"
	"npool/gopool"
	"time"
	"zhubeat/beat/job"
)

type ClientArgs struct {
	network   string
	addresses []string // ["127.0.0.1:7070","127.0.0.1:7070","127.0.0.1:7070"]
	timeout   time.Duration
	ttl       time.Duration
	// retry
}

func NewClientArgs(network string, addresses []string, timeout, ttl time.Duration) *ClientArgs {
	return &ClientArgs{
		network:   network,
		addresses: addresses,
		timeout:   timeout,
		ttl:       ttl,
	}
}

type ConnPool struct {
	address []string
	network string // tcp
	*gopool.ChanConnPool
}

func NewConnPool(network string, address []string) (*ConnPool, error) {
	pool, err := gopool.NewChanConnPool(&gopool.ConnPoolReq{
		Addrs: address,
		Create: func(addr string, timeout time.Duration) (interface{}, error) {
			cli, err := net.DialTimeout(network, addr, timeout)
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
	return &ConnPool{
		address:      address,
		network:      network,
		ChanConnPool: pool,
	}, nil
}

type Client struct {
	*ConnPool
	name    string
	timeout time.Duration
	ttl     time.Duration
}

func NewClient(name string, timeout, ttl time.Duration, pool *ConnPool) *Client {
	return &Client{
		name:     name,
		ConnPool: pool,
		timeout:  timeout,
		ttl:      ttl,
	}
}

func (c *Client) getConn() (conn net.Conn, err error) {
	var ok bool
	cli, err := c.Get()
	if err != nil {
		return
	}
	if conn, ok = cli.(net.Conn); !ok {
		err = errors.New("cli is nil.")
	}
	return
}

func (c *Client) putConn(conn net.Conn, err *error) {
	safe := false
	if *err == nil {
		safe = true
	}
	c.Put(conn, safe)
}

func (c *Client) reGetConn(src net.Conn, srcErr *error) (conn net.Conn, err error) {
	c.putConn(src, srcErr)
	return c.getConn()
}

// 如果失败下一次则重新在连接池里取conn
func (c *Client) Publish(batch job.BatchJob) (err error) {
	var (
		conn       net.Conn
		needReconn bool
		firstConn  bool
	)

	defer func() {
		c.putConn(conn, &err)
	}()

	firstConn = true

	for _, message := range batch {
		time.Sleep(c.ttl)

		fmt.Println("c.GetHealthy():", c.GetHealthy())
		// fmt.Println("firstConn is ", firstConn)
		// fmt.Println("needReconn is ", needReconn)

		if firstConn {
			conn, err = c.getConn()
			if err != nil {
				fmt.Println("1111 store message: ", message)
				firstConn = false
				needReconn = true
				continue
			}
			firstConn = false
		}
		if needReconn {
			conn, err = c.reGetConn(conn, &err)
			if err != nil {
				fmt.Println("222 store message: ", message)
				continue
			}
		}
		_, err := conn.Write([]byte(string(message)))
		if err != nil {
			// store
			fmt.Println("333 store message: ", message)
			needReconn = true
		} else {
			needReconn = false
		}
	}

	return err
}

func (c *Client) Close() error {
	// 应该return conn
	// 所以c 里应该有当前的conn
	return nil
}
