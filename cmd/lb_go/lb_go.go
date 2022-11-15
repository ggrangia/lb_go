package lb_go

import (
	"fmt"
	"log"
	"math/rand"
	"net/http/httputil"
	"net/url"
)

type Lb struct {
	Backends []backend
	Selector Selector
}

type Selector interface {
	Select([]backend) backend
}
type RandomSelection struct {
	Seed int64
}

func (rs *RandomSelection) Select(backends []backend) backend {
	rand.Seed(rs.Seed)
	return backends[rand.Intn(len(backends))]
}

type backend struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

func NewBackend(myurl string) backend {
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
