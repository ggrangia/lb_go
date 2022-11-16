package lb_go

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func setupBackends(t *testing.T) (func(t *testing.T), []Backend) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	backendServer2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy2")
	}))

	backends := []Backend{
		NewBackend(backendServer.URL),
		NewBackend(backendServer2.URL),
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
	expected := Backend{
		Addr:  wantUrl,
		Proxy: httputil.NewSingleHostReverseProxy(rpURL),
	}
	result := NewBackend(wantUrl)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func TestRandomSelectionFull(t *testing.T) {
	teardown, backends := setupBackends(t)
	defer teardown(t)
	rs := RandomSelection{
		Seed: time.Now().UTC().UnixNano(),
	}

	lb := Lb{
		Backends: backends,
		Selector: &rs,
	}

	frontendProxy := httptest.NewServer(http.HandlerFunc(lb.Selector.Select(lb.Backends).Proxy.ServeHTTP))
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
