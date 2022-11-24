package roundrobin

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ggrangia/lb_go/pkg/backend"
)

type RoundRobin struct {
	mutex    sync.RWMutex
	Counter  int
	backends []backend.Backend
}

func (rr *RoundRobin) AddBackend(b backend.Backend) {
	rr.backends = append(rr.backends, b)
}

func (rr *RoundRobin) Select() backend.Backend {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	attempts := 1
	for {
		b := rr.backends[rr.Counter]
		rr.Counter = (rr.Counter + 1) % len(rr.backends)
		if b.Alive {
			return b
		}
		fmt.Println()
		if attempts >= len(rr.backends) {
			panic("AHHHHH, none is alive!!!!")
		}
		attempts++
	}
}

func (rr *RoundRobin) nextServer() backend.Backend {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	attempts := 1
	for {
		b := rr.backends[rr.Counter]
		rr.Counter = (rr.Counter + 1) % len(rr.backends)
		if b.Alive {
			return b
		}
		fmt.Println()
		if attempts >= len(rr.backends) {
			panic("AHHHHH, none is alive!!!!")
		}
		attempts++
	}
}

func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr.nextServer().Proxy.ServeHTTP(w, r)
}

func New() *RoundRobin {
	return &RoundRobin{}
}

func NewWithBackends(backends []backend.Backend) *RoundRobin {
	return &RoundRobin{
		backends: backends,
	}
}
