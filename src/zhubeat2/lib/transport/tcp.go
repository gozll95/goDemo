package transport

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func NetDialer(timeout time.Duration) Dialer {
	return DialerFunc(func(network, address string) (net.Conn, error) {
		switch network {
		case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
		default:
			return nil, fmt.Errorf("unsupported network type %v", network)
		}

		addresses := strings.Split(address, string(DELIM))

		// dial via host IP by randomized iteration of known IPs
		dialer := &net.Dialer{Timeout: timeout}
		return DialWith(dialer, network, addresses)
	})
}
