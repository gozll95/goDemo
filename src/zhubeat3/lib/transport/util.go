package transport

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func fullAddress(host string, defaultPort int) string {
	if _, _, err := net.SplitHostPort(host); err == nil {
		return host
	}

	idx := strings.Index(host, ":")
	if idx >= 0 {
		// IPv6 address detected
		return fmt.Sprintf("[%v]:%v", host, defaultPort)
	}
	return fmt.Sprintf("%v:%v", host, defaultPort)
}

// DialWith randomly dials one of a number of addresses with a given dialer.
//
// Use this to select and dial one IP being known for one host name.
func DialWith(
	dialer Dialer,
	network string, // ["xx:xx","xx:xx"] string
	addresses []string, // ["xx:xx","xx:xx"]
) (c net.Conn, err error) {
	//fmt.Println("addresses is ", addresses)
	switch len(addresses) {
	case 0:
		return nil, fmt.Errorf("no route to host %v", addresses)
	case 1:
		return dialer.Dial(network, addresses[0])
	}

	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(addresses)) {
		c, err = dialer.Dial(network, addresses[i])

		if err == nil && c != nil {
			//fmt.Println("addresses[i] is ", "i is", addresses[i], i)
			return c, err
		}
	}

	if err == nil {
		err = fmt.Errorf("unable to connect to '%v'", addresses)
	}
	return nil, err
}
