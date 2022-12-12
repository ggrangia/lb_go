package randomselection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ggrangia/lb_go/test"
	"github.com/google/go-cmp/cmp"
)

func TestNoBackends(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	seq_len := 10

	rs := New(seed)

	for i := 0; i < seq_len; i++ {
		got, err := rs.nextServer()
		if err == nil {
			t.Fatalf("No error connecting. Got: %v", got)
		}
		if !cmp.Equal(err.Error(), ErrNoServer.Error()) {
			t.Errorf("Expected output %v got %v", ErrNoServer.Error(), err.Error())
		}
	}
}

func TestRandomSelection(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)

	rs := NewWithBackends(seed, backends)

	// Backends alive is false by default
	for i := 0; i < len(backends); i++ {
		rs.AddBackend(backends[i])
	}

	source := rand.NewSource(seed)
	generator := rand.New(source)
	seq_len := 10
	for i := 0; i < seq_len; i++ {
		got, err := rs.nextServer()
		if err != nil {
			t.Fatalf("Error connecting. Got: %v", err)
		}

		b := rs.Backends[generator.Intn(len(rs.Backends))]

		if !cmp.Equal(b.Url, got.Url) {
			t.Errorf("Got %v, want %v", b.Url, got.Url)
		}
	}
}
