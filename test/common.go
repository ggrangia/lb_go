package test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggrangia/lb_go/pkg/backend"
)

func setupBackends(t *testing.T, n int) (func(t *testing.T), []backend.Backend) {

	backends := make([]backend.Backend, n)
	servers := make([]*httptest.Server, n)
	for i := 0; i < n; i++ {
		str := fmt.Sprintf("this call was relayed by the reverse proxy%d\n", i)
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, str)
		}))
		backends[i] = backend.NewBackend(servers[i].URL)
	}

	teardown := func(t *testing.T) {
		for i := 0; i < len(backends); i++ {
			defer servers[i].Close()

		}
	}
	fmt.Println(servers)
	fmt.Println(backends)
	return teardown, backends

}

func getTest(url string) {
	// GET test
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
}
