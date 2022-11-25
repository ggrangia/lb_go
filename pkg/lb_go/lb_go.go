package lb_go

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ggrangia/lb_go/pkg/healthcheck"
	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
	"github.com/ggrangia/lb_go/pkg/lb_go/selection"
)

type Lb struct {
	Backends     []*backend.Backend
	Selector     selection.Selector
	health_timer int
}

func NewLb(backends []*backend.Backend, selector selection.Selector) *Lb {

	return &Lb{
		Backends:     backends,
		Selector:     selector,
		health_timer: 5, // default value
	}
}

func (lb *Lb) SetHealthcheckTimer(timer int) {
	lb.health_timer = timer
}

func (lb *Lb) Start() {
	lb_proxy := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: http.HandlerFunc(lb.Selector.ServeHTTP),
	}

	go lb.runHealthchecks()

	if err := lb_proxy.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (lb *Lb) runHealthchecks() {
	ticker := time.NewTicker(time.Second * time.Duration(lb.health_timer))
	for range ticker.C {
		lb.healthchecks()
	}
}

func (lb *Lb) healthchecks() {
	for _, b := range lb.Selector.Backends {
		alive := healthcheck.IsAliveTCP(b.Url)
		fmt.Printf("%v is %v, it becomes %v\n", b.Addr, b.Alive, alive)
		b.Alive = alive
	}
}
