package selection

import (
	"fmt"
)

type RoundRobin struct {
	Counter int
}

func (rr *RoundRobin) Select(l int) int {
	b := rr.Counter % l
	fmt.Printf("Current counter %d\n", rr.Counter)
	rr.Counter += 1
	return b
}

func NewRoundRobin() Selector {
	return &RoundRobin{}
}
