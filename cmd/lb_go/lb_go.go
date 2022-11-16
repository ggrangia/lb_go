package lb_go

import (
	"fmt"
	"log"
	"math/rand"
	"net/http/httputil"
	"net/url"
)

type Lb struct {
	Backends []Backend
	Selector Selector
}

type Selector interface {
	Select([]Backend) Backend
}
type RandomSelection struct {
	Seed int64
}

func (rs *RandomSelection) Select(backends []Backend) Backend {
	rand.Seed(rs.Seed)
	return backends[rand.Intn(len(backends))]
}

type Backend struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

func NewBackend(myurl string) Backend {
	rpURL, err := url.Parse(myurl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rpURL)
	return Backend{
		Addr:  myurl,
		Proxy: httputil.NewSingleHostReverseProxy(rpURL),
	}
}
