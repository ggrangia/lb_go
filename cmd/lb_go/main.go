package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/ggrangia/lb_go/pkg/backend"
	"github.com/ggrangia/lb_go/pkg/lb_go"
	"github.com/ggrangia/lb_go/pkg/selection"
	"github.com/ggrangia/lb_go/pkg/selection/randomselection"
	"github.com/ggrangia/lb_go/pkg/selection/roundrobin"
)

func main() {

	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	defer backendServer.Close()

	backendServer2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy2")
	}))
	defer backendServer2.Close()
	backendServer3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy3")
	}))
	defer backendServer3.Close()
	backends := []backend.Backend{
		backend.NewBackend(backendServer.URL),
		backend.NewBackend(backendServer2.URL),
		backend.NewBackend(backendServer3.URL),
	}

	algo := "roundrobin"
	var selector selection.Selector
	switch algo {
	case "roundrobin":
		selector = roundrobin.NewRoundRobin()
	case "randomselection":
		selector = randomselection.NewRandomSelection(time.Now().UTC().UnixNano())
	default:
		log.Fatal("Selection algorithm unknown")
	}
	lb := lb_go.NewLb(backends, selector)
	lb.Start()
}
