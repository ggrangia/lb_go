package integration_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ggrangia/lb_go/pkg/lb_go"
	"github.com/ggrangia/lb_go/pkg/lb_go/selection/roundrobin"
	"github.com/ggrangia/lb_go/test"
)

func TestE2eRoundRobin(t *testing.T) {
	var wg sync.WaitGroup
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)

	selector := roundrobin.New()

	lb := lb_go.Lb{
		Backends: backends,
		Selector: selector,
	}
	frontendProxy := httptest.NewServer(http.HandlerFunc(lb.Serve))
	defer frontendProxy.Close()

	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			test.GetTest(frontendProxy.URL)
			wg.Done()
		}()
	}
	wg.Wait()

	t.Errorf("FIXME: To be Completed")
}
