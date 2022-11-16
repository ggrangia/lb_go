package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggrangia/lb_go/pkg/backend"
)

func setupBackends(t *testing.T) (func(t *testing.T), []backend.Backend) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	backendServer2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy2")
	}))

	backends := []backend.Backend{
		backend.NewBackend(backendServer.URL),
		backend.NewBackend(backendServer2.URL),
	}
	teardown := func(t *testing.T) {
		defer backendServer.Close()
		defer backendServer2.Close()
	}
	return teardown, backends

}
