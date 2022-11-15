package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
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

	backends := []backend{
		newBackend(backendServer.URL),
		newBackend(backendServer2.URL),
	}

	rs := RandomSelection{
		seed: time.Now().UTC().UnixNano(),
	}

	lb := Lb{
		backends: backends,
		selector: &rs,
	}
	/*
		lb_proxy := http.Server{
			Addr:    fmt.Sprintf(":%d", 8080),
			Handler: http.HandlerFunc(lb.random_selection),
		}
		if err := lb_proxy.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	*/
	frontendProxy := httptest.NewServer(http.HandlerFunc(lb.selector.Select(lb.backends).Proxy.ServeHTTP))
	defer frontendProxy.Close()

	// GET test
	resp, err := http.Get(frontendProxy.URL)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)

}
