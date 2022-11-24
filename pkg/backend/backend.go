package backend

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	Addr  string
	Proxy *httputil.ReverseProxy
	Url   *url.URL
	Alive bool
	http.Handler
}

func NewBackend(myurl string) Backend {
	rpURL, err := url.Parse(myurl)
	if err != nil {
		log.Fatal(err)
	}

	return Backend{
		Addr:  myurl,
		Url:   rpURL,
		Proxy: httputil.NewSingleHostReverseProxy(rpURL),
		Alive: false,
	}
}
