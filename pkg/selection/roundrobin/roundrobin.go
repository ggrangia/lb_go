package roundrobin

import (
	"sync/atomic"
)

type RoundRobin struct {
	Counter uint64
}

func (rr *RoundRobin) Select(l int) int {
	// First return is 1
	return int(atomic.AddUint64(&rr.Counter, uint64(1)) % uint64(l))
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}
