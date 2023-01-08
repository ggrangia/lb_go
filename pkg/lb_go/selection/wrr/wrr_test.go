package wrr

import (
	"testing"

	"github.com/ggrangia/lb_go/test"
	"github.com/google/go-cmp/cmp"
)

func TestNextServerAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 6)
	weights := [6]int{7, 9, 5, 3, 11, 4}
	defer teardown(t)
	rr := New()

	for i := 0; i < len(backends); i++ {
		// Set all backends alive
		backends[i].Alive = true
		rr.AddWeightedBackend(backends[i], weights[i])
	}

	tests := []struct {
		want int
	}{
		{want: 11},
		{want: 9},
		{want: 7},
		{want: 11},
		{want: 5},
		{want: 9},
		{want: 4},
		{want: 11},
		{want: 7},
		{want: 3},
	}

	for i, tc := range tests {
		got, err := rr.nextServer()
		if err != nil {
			t.Fatalf("Got error %v", err)
		}
		t.Logf("Iter: %d, Got %v %v", i, got.weight, got.Url.Host)

		if !cmp.Equal(got.weight, tc.want) {
			t.Errorf("Expected output %v got %v", tc.want, got.weight)
		}
	}
}

func TestNextServerNotAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 6)
	weights := [6]int{7, 9, 5, 3, 11, 4}
	defer teardown(t)
	rr := New()

	for i := 0; i < len(backends); i++ {
		rr.AddWeightedBackend(backends[i], weights[i])
	}

	tests := []struct {
		want int
	}{
		{want: 11},
		{want: 9},
		{want: 7},
		{want: 11},
		{want: 5},
		{want: 9},
		{want: 4},
		{want: 11},
		{want: 7},
		{want: 3},
	}

	for _, tc := range tests {
		got, err := rr.nextServer()
		if err == nil {
			t.Fatalf("No error connecting to %v. Got: %v", tc.want, got)
		}

		if !cmp.Equal(err.Error(), ErrNoServer.Error()) {
			t.Errorf("Expected output %v got %v", ErrNoServer.Error(), err.Error())
		}
	}
}
