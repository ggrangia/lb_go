package randomselection

import (
	"math/rand"
)

type RandomSelection struct {
	generator rand.Rand
}

func (rs *RandomSelection) Select(l int) int {
	return rs.generator.Intn(l)
}

func NewRandomSelection(seed int64) *RandomSelection {
	source := rand.NewSource(seed)
	generator := rand.New(source)
	return &RandomSelection{
		generator: *generator,
	}
}