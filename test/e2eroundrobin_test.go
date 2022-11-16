package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggrangia/lb_go/pkg/lb_go"
	"github.com/ggrangia/lb_go/pkg/selection"
)

func TestE2eRoundRobin(t *testing.T) {
	teardown, backends := setupBackends(t, 3)
	defer teardown(t)

	selector := selection.RoundRobin{}

	lb := lb_go.Lb{
		Backends: backends,
		Selector: &selector,
	}
	frontendProxy := httptest.NewServer(http.HandlerFunc(lb.Serve))
	defer frontendProxy.Close()

	getTest(frontendProxy.URL)
	getTest(frontendProxy.URL)

	t.Errorf("FIXME: To be Completed")
}
