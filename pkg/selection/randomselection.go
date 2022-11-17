package selection

import (
	"math/rand"
)

type RandomSelection struct {
	Seed int64
}

func (rs *RandomSelection) Select(l int) int {
	rand.Seed(rs.Seed)
	return rand.Intn(l)
}
