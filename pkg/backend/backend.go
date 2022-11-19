package backend

import (
	"log"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

func NewBackend(myurl string) Backend {
	rpURL, err := url.Parse(myurl)
	if err != nil {
		log.Fatal(err)
	}

	return Backend{
		Addr:  myurl,
		Proxy: httputil.NewSingleHostReverseProxy(rpURL),
	}
}
