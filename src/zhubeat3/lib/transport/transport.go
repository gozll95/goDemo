package transport

import (
	"errors"
	"net"
	"time"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type DialerFunc func(network, address string) (net.Conn, error)

var (
	ErrNotConnected = errors.New("client is not connected")
)

func (d DialerFunc) Dial(network, address string) (net.Conn, error) {
	return d(network, address)
}

func Dial(network, address string, timeout time.Duration) (net.Conn, error) {
	d, err := MakeDialer(timeout)
	if err != nil {
		return nil, err
	}
	return d.Dial(network, address)
}
