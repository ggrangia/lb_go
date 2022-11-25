package backend_test

import (
	"log"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
	"github.com/google/go-cmp/cmp"
)

func TestNewBackend(t *testing.T) {
	wantUrl := "http://127.0.0.1:11111"
	rpURL, err := url.Parse(wantUrl)
	if err != nil {
		log.Fatal(err)
	}
	expected := backend.Backend{
		Addr:  wantUrl,
		Proxy: httputil.NewSingleHostReverseProxy(rpURL),
	}
	result := backend.NewBackend(wantUrl)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Got backend %v want %v, diff: %v", result, expected, diff)
	}
}
