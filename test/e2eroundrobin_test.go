package test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ggrangia/lb_go/pkg/lb_go"
	"github.com/ggrangia/lb_go/pkg/selection/roundrobin"
)

func TestE2eRoundRobin(t *testing.T) {
	var wg sync.WaitGroup
	teardown, backends := setupBackends(t, 3)
	defer teardown(t)

	selector := roundrobin.RoundRobin{}

	lb := lb_go.Lb{
		Backends: backends,
		Selector: &selector,
	}
	frontendProxy := httptest.NewServer(http.HandlerFunc(lb.Serve))
	defer frontendProxy.Close()

	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			getTest(frontendProxy.URL)
			wg.Done()
		}()
	}
	wg.Wait()

	//t.Errorf("FIXME: To be Completed")
}
