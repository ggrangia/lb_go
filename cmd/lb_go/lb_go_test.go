package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setupBackends(t *testing.T) (func(t *testing.T), []backend) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	backendServer2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy2")
	}))

	backends := []backend{
		newBackend(backendServer.URL),
		newBackend(backendServer2.URL),
	}
	teardown := func(t *testing.T) {
		defer backendServer.Close()
		defer backendServer2.Close()
	}
	return teardown, backends

}

func TestNewBackend(t *testing.T) {
	wantUrl := "http://127.0.0.1:11111"
	rpURL, err := url.Parse(wantUrl)
	if err != nil {
		log.Fatal(err)
	}
	expected := backend{
		Addr:  wantUrl,
		Proxy: httputil.NewSingleHostReverseProxy(rpURL),
	}
	result := newBackend(wantUrl)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func TestRandomSelectionFull(t *testing.T) {
	teardown, backends := setupBackends(t)
	defer teardown(t)
	lb := Lb{backends}

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
	t.Fatalf("To be completed")
}
