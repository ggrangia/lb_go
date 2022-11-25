package lb_go

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ggrangia/lb_go/pkg/healthcheck"
	"github.com/ggrangia/lb_go/pkg/lb_go/selection"
)

type Lb struct {
	Selector       selection.Selector
	health_service *healthcheck.Healthchecker
}

func NewLb(selector selection.Selector, hs *healthcheck.Healthchecker) *Lb {

	return &Lb{
		Selector:       selector,
		health_service: hs,
	}
}

func (lb *Lb) Start() {
	lb_proxy := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: http.HandlerFunc(lb.Selector.ServeHTTP),
	}

	go lb.health_service.RunHealthchecks()
	if err := lb_proxy.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
