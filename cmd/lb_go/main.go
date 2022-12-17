package main

import (
	"github.com/ggrangia/lb_go/pkg/cmd"
)

func main() {

	cmd.Execute()
	/*
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
		backends := []*backend.Backend{
			backend.NewBackend(backendServer.URL),
			backend.NewBackend(backendServer2.URL),
			backend.NewBackend(backendServer3.URL),
		}

		algo := "roundrobin"
		var selector selection.Selector
		switch algo {
		case "roundrobin":
			selector = roundrobin.NewWithBackends(backends)
		case "randomselection":
			selector = randomselection.NewWithBackends(time.Now().UTC().UnixNano(), backends)
		default:
			log.Fatalf("Unknown selection algorithm: %v", algo)
		}

		hc := healthcheck.New(selector, 5)

		lb := lb_go.NewLb(selector, hc)
		lb.Start()
	*/
}
