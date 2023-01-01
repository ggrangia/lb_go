package wrr

import (
	"sync"

	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
)

type Wrr struct {
	mutex    sync.RWMutex
	Backends []*backend.Backend
}

func New() *Wrr {
	return &Wrr{}
}

func NewWithBackends(backends []*backend.Backend) *Wrr {
	return &Wrr{
		Backends: backends,
	}
}
