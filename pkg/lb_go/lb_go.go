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
	URL            string
	port           int
	proxy          *http.Server
}

func NewLb(selector selection.Selector, hs *healthcheck.Healthchecker, port int) *Lb {

	return &Lb{
		Selector:       selector,
		health_service: hs,
		port:           port,
	}
}

func (lb *Lb) Start() {
	lb_proxy := &http.Server{
		Addr:    fmt.Sprintf(":%d", lb.port),
		Handler: http.HandlerFunc(lb.Selector.ServeHTTP),
	}

	lb.URL = lb_proxy.Addr
	lb.proxy = lb_proxy

	go lb.health_service.RunHealthchecks()
	if err := lb_proxy.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
