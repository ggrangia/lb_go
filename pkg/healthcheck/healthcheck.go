package healthcheck

import (
	"log"
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

func New(s selection.Selector, i time.Duration) *Healthchecker {
	log.Printf("Starting healthckecks with timer %d\n", i)
	return &Healthchecker{
		Selector: s,
		interval: i,
	}
}

func (hs *Healthchecker) IsAliveTCP(url *url.URL) bool {
	timeout := time.Second * 5
	conn, err := net.DialTimeout("tcp", url.Host, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}

func (hs *Healthchecker) RunHealthchecks() {
	// Call healthchecks immidiately
	hs.healthchecks()
	ticker := time.NewTicker(time.Second * hs.interval)
	for range ticker.C {
		hs.healthchecks()
	}
}

func (hs *Healthchecker) healthchecks() {
	for _, b := range hs.Selector.GetBackends() {
		alive := hs.IsAliveTCP(b.Url)
		log.Printf("%v is %v, it becomes %v\n", b.Addr, b.Alive, alive)
		b.Alive = alive
	}
}

func (hs *Healthchecker) SetHealthcheckTimer(interval int) {
	hs.interval = time.Duration(interval)
}
