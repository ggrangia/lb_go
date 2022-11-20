package healthcheck_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ggrangia/lb_go/pkg/healthcheck"
)

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
	res = healthcheck.IsAliveTCP(myurl)
	if !res {
		t.Errorf("Expected %v to be alive", myurl)
	}

	// Offline test - res must be false
	offlineStr := "localhost:8080"
	offlineUrl, err := url.Parse(offlineStr)
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}
	res = healthcheck.IsAliveTCP(offlineUrl)
	if res {
		t.Errorf("Expected %v to be dead", offlineUrl)
	}

}
