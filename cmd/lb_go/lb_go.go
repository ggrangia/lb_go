package main

import (
	"fmt"
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
	selector Selector
}

type Selector interface {
	Select ([]backend) backend
} 
type RandomSelection struct {
	seed int64
}

func (rs *RandomSelection) Selector(backends []backend) {
	rand.Seed(rs.seed)
	return backends[rand.Intn(len(lb.backends))]
}

type backend struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

func newBackend(myurl string) backend {
	rpURL, err := url.Parse(myurl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rpURL)
	return backend{
		Addr:  myurl,
		Proxy: httputil.NewSingleHostReverseProxy(rpURL),
	}
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
		newBackend(backendServer.URL),
		newBackend(backendServer2.URL),
	}

	lb := Lb{
		backends: backends,
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
