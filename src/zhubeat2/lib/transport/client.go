package transport

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const DELIM = ','

type Hosts string

func (h Hosts) String() string {
	return string(h)
}
func (h Hosts) Slice() []string {
	return strings.Split(h.String(), string(DELIM))
}

func ConvertToHosts(hosts []string) Hosts {
	return Hosts(strings.Join(hosts, string(DELIM)))
}

type Client struct {
	dialer  Dialer
	network string
	hosts   Hosts

	conn  net.Conn
	mutex sync.Mutex
}

func MakeDialer(timeout time.Duration) (Dialer, error) {
	var err error
	dialer := NetDialer(timeout)
	if err != nil {
		return nil, err
	}
	return dialer, nil
}

func NewClient(network string, hosts []string, timeout time.Duration) (*Client, error) {
	// do some sanity checks regarding network and Config matching +
	// address being parseable
	switch network {
	case "tcp", "tcp4", "tcp6":
	case "udp", "udp4", "udp6":
	default:
		return nil, fmt.Errorf("unsupported network type %v", network)
	}

	dialer, err := MakeDialer(timeout)
	if err != nil {
		return nil, err
	}

	return NewClientWithDialer(dialer, network, hosts)
}

func NewClientWithDialer(d Dialer, network string, hosts []string) (*Client, error) {
	// check address being parseable
	for _, host := range hosts {
		_, _, err := net.SplitHostPort(host)
		if err != nil {
			return nil, err
		}
	}

	client := &Client{
		dialer:  d,
		network: network,
		hosts:   ConvertToHosts(hosts),
	}
	return client, nil
}

// Connect() will return a new net.Conn every time
func (c *Client) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}

	conn, err := c.dialer.Dial(c.network, c.hosts.String())
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) IsConnected() bool {
	c.mutex.Lock()
	b := c.conn != nil
	c.mutex.Unlock()
	return b
}

func (c *Client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		fmt.Println("close conn", c.conn)
		c.conn = nil
		return err
	}
	return nil
}

func (c *Client) getConn() net.Conn {
	c.mutex.Lock()
	conn := c.conn
	c.mutex.Unlock()
	return conn
}

func (c *Client) Read(b []byte) (int, error) {
	conn := c.getConn()
	if conn == nil {
		return 0, ErrNotConnected
	}

	n, err := conn.Read(b)
	return n, c.handleError(err)
}

func (c *Client) Write(b []byte) (int, error) {
	conn := c.getConn()
	if conn == nil {
		return 0, ErrNotConnected
	}

	n, err := c.conn.Write(b)
	return n, c.handleError(err)
}

func (c *Client) LocalAddr() net.Addr {
	conn := c.getConn()
	if conn != nil {
		return c.conn.LocalAddr()
	}
	return nil
}

func (c *Client) RemoteAddr() net.Addr {
	conn := c.getConn()
	if conn != nil {
		return c.conn.LocalAddr()
	}
	return nil
}

func (c *Client) Host() string {
	return c.hosts.String()
}

func (c *Client) SetDeadline(t time.Time) error {
	conn := c.getConn()
	if conn == nil {
		return ErrNotConnected
	}

	err := conn.SetDeadline(t)
	return c.handleError(err)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	conn := c.getConn()
	if conn == nil {
		return ErrNotConnected
	}

	err := conn.SetReadDeadline(t)
	return c.handleError(err)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	conn := c.getConn()
	if conn == nil {
		return ErrNotConnected
	}

	err := conn.SetWriteDeadline(t)
	return c.handleError(err)
}

func (c *Client) handleError(err error) error {
	if err != nil {
		if nerr, ok := err.(net.Error); !(ok && (nerr.Temporary() || nerr.Timeout())) {
			_ = c.Close()
		}
	}
	return err
}
