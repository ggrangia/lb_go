package lb_go

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ggrangia/lb_go/pkg/backend"
	"github.com/ggrangia/lb_go/pkg/healthcheck"
	"github.com/ggrangia/lb_go/pkg/selection"
)

type Lb struct {
	Backends []backend.Backend
	Selector selection.Selector
}

func (lb *Lb) Serve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called the proxy")
	b := lb.Selector.Select(len(lb.Backends))

	lb.Backends[b].Proxy.ServeHTTP(w, r)
}

// FIXME: init the lb server
// FIXME: make interval variable
func (lb *Lb) Start() {
	go lb.runHealthchecks(10)
}

func (lb *Lb) runHealthchecks(seconds int) {
	ticker := time.NewTicker(time.Second * time.Duration(seconds))
	for range ticker.C {
		lb.healthchecks()
	}
}

func (lb *Lb) healthchecks() {
	for _, b := range lb.Backends {
		alive := healthcheck.IsAliveTCP(b.Url)
		b.Alive = alive
	}
}
