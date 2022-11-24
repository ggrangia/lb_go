package roundrobin_test

import (
	"testing"

	"github.com/ggrangia/lb_go/pkg/lb_go/selection/roundrobin"
	"github.com/ggrangia/lb_go/test"
	"github.com/google/go-cmp/cmp"
)

func TestRoundRobinAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)
	rr := roundrobin.New()

	for i := 0; i < len(backends); i++ {
		// Set all backends alive
		backends[i].Alive = true
		rr.AddBackend(backends[i])
	}

	tests := []struct {
		want    string
		counter int
	}{
		{want: backends[0].Url.Host, counter: 1},
		{want: backends[1].Url.Host, counter: 2},
		{want: backends[2].Url.Host, counter: 0},
		{want: backends[0].Url.Host, counter: 1},
	}

	for _, tc := range tests {
		got := rr.Select()
		if !cmp.Equal(got.Url.Host, tc.want) {
			t.Errorf("Expected output %v got %v", tc.want, got.Url.Host)
		}
		if !cmp.Equal(rr.Counter, tc.counter) {
			t.Errorf("Expected counter %d got %d", tc.counter, rr.Counter)
		}
	}
}

func TestRoundRobinNotAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)
	rr := roundrobin.New()

	// Backends alive is false by default
	for i := 0; i < len(backends); i++ {
		rr.AddBackend(backends[i])
	}

	t.Fatal("TO IMPLEMENT")
}
