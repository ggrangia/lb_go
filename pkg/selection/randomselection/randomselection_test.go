package randomselection_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ggrangia/lb_go/pkg/selection/randomselection"
	"github.com/google/go-cmp/cmp"
)

func TestRandomSelection(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	seq_len := 10
	max_value := 15

	got := make([]int, seq_len)
	want := make([]int, seq_len)

	source := rand.NewSource(seed)
	generator := rand.New(source)

	rs := randomselection.NewRandomSelection(seed)

	for i := 0; i < seq_len; i++ {
		got[i] = rs.Select(max_value)
		want[i] = generator.Intn(max_value)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Got %d, want %d", got, want)
	}
}
