package healthcheck

import (
	"net"
	"net/url"
	"time"

	"github.com/ggrangia/lb_go/pkg/lb_go/selection"
)

// Take the Selector
type Healthchecker struct {
	Selector selection.Selector
	interval time.Duration
}

func IsAliveTCP(url *url.URL) bool {
	timeout := time.Second * 5
	conn, err := net.DialTimeout("tcp", url.Host, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}
