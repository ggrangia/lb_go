package selection

import (
	"math/rand"

	"github.com/ggrangia/lb_go/pkg/backend"
)

type RandomSelection struct {
	Seed int64
}

func (rs *RandomSelection) Select(backends []backend.Backend) backend.Backend {
	rand.Seed(rs.Seed)
	return backends[rand.Intn(len(backends))]
}
