package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	lb "github.com/ggrangia/lb_go/cmd/lb_go"
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
	backends := []lb.Backend{
		lb.NewBackend(backendServer.URL),
		lb.NewBackend(backendServer2.URL),
	}

	// FIXME: fetch Selector

	lb := lb.Lb{
		Backends: backends,
		Selector: &rs,
	}

	lb_proxy := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: http.HandlerFunc(lb.Selector.Select(lb.Backends).Proxy.ServeHTTP),
	}
	if err := lb_proxy.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
