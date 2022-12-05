package healthcheck_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ggrangia/lb_go/pkg/healthcheck"
	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
)

type mockSelector struct {
	B []*backend.Backend
}

func (fs *mockSelector) GetBackends() []*backend.Backend {
	return fs.B
}

func (fs *mockSelector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.B[0].ServeHTTP(w, r)
}

func TestIsAliveTCP(t *testing.T) {

	var res bool

	// Online test - res must be true
	str := "Alive Backend"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, str)
	}))
	myurl, err := url.Parse(s.URL)
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	selec := mockSelector{B: []*backend.Backend{backend.NewBackend(s.URL)}}

	hc := healthcheck.New(&selec, 5)
	res = hc.IsAliveTCP(myurl)
	if !res {
		t.Errorf("Expected %v to be alive", myurl)
	}

	// Offline test - res must be false
	offlineStr := "localhost:8080"
	offlineUrl, err := url.Parse(offlineStr)
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}
	res = hc.IsAliveTCP(offlineUrl)
	if res {
		t.Errorf("Expected %v to be dead", offlineUrl)
	}

}
