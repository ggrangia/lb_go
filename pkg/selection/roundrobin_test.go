package selection_test

import (
	"testing"

	"github.com/ggrangia/lb_go/pkg/selection"
	"github.com/google/go-cmp/cmp"
)

func TestRoundRobin(t *testing.T) {
	tests := []struct {
		input   int
		want    int
		counter int
	}{
		{input: 3, want: 0, counter: 0},
		{input: 3, want: 1, counter: 1},
		{input: 3, want: 2, counter: 2},
		{input: 3, want: 0, counter: 3},
		{input: 3, want: 1, counter: 4},
	}
	rr := selection.RoundRobin{}

	for c, tc := range tests {
		got := rr.Select(tc.input)
		if !cmp.Equal(got, tc.want) {
			t.Errorf("Expected output %d got %d", tc.want, got)
		}
		if !cmp.Equal(c, tc.counter) {
			t.Errorf("Expected counter %d got %d", tc.counter, c)
		}
	}
}
