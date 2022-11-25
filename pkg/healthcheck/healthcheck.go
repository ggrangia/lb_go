package healthcheck

import (
	"fmt"
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

func (hs *Healthchecker) isAliveTCP(url *url.URL) bool {
	timeout := time.Second * 5
	conn, err := net.DialTimeout("tcp", url.Host, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}

func (hs *Healthchecker) RunHealthchecks() {
	ticker := time.NewTicker(time.Second * hs.interval)
	for range ticker.C {
		hs.healthchecks()
	}
}

func (hs *Healthchecker) healthchecks() {
	for _, b := range hs.Selector.GetBackends() {
		alive := hs.isAliveTCP(b.Url)
		fmt.Printf("%v is %v, it becomes %v\n", b.Addr, b.Alive, alive)
		b.Alive = alive
	}
}

func (hs *Healthchecker) SetHealthcheckTimer(interval int) {
	hs.interval = time.Duration(interval)
}
