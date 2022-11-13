package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"
)

type Lb struct {
	backends []backend
}

type backend struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func NewBackend(myurl string) backend {
	rpURL, err := url.Parse(myurl)
	if err != nil {
		log.Fatal(err)
	}
	return backend{
		addr:  myurl,
		proxy: httputil.NewSingleHostReverseProxy(rpURL),
	}
}

func (lb *Lb) random_selection(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	lb.backends[rand.Intn(len(lb.backends))].proxy.ServeHTTP(w, r)
}

/*
	func (lb *Lb) lb_algo(w http.ResponseWriter, r *http.Request) {
		lb.random_selection().proxy.ServeHTTP(w, r)
	}
*/
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
		NewBackend(backendServer.URL),
		NewBackend(backendServer2.URL),
	}

	lb := Lb{
		backends: backends,
	}

	proxyHandler := func(r http.ResponseWriter, w *http.Request) {
		lb.random_selection(r, w)
	}
	frontendProxy := httptest.NewServer(http.HandlerFunc(proxyHandler))
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
