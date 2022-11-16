package selection

import (
	"fmt"

	"github.com/ggrangia/lb_go/pkg/backend"
)

type RoundRobin struct {
	Counter int
}

func (rr *RoundRobin) Select(backends []backend.Backend) backend.Backend {
	b := backends[rr.Counter%len(backends)]
	fmt.Printf("Current counter %d\n", rr.Counter)
	rr.Counter += 1
	return b
}
