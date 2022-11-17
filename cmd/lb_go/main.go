package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/ggrangia/lb_go/pkg/backend"
	lb "github.com/ggrangia/lb_go/pkg/lb_go"
	"github.com/ggrangia/lb_go/pkg/selection"
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

	// FIXME: fetch Selector
	//rs := selection.RoundRobin{}
	rs := selection.NewRandomSelection(time.Now().UTC().UnixNano())
	lb := lb.Lb{
		Backends: backends,
		Selector: rs,
	}

	lb_proxy := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: http.HandlerFunc(lb.Serve),
	}

	if err := lb_proxy.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
