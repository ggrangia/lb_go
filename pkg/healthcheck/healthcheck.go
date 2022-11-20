package healthcheck

import (
	"net"
	"net/url"
	"time"
)

func IsAliveTCP(url *url.URL) bool {
	timeout := time.Second * 5
	conn, err := net.DialTimeout("tcp", url.Host, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}
