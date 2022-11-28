package roundrobin

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
)

var ErrNoServer = errors.New("No available servers")

type RoundRobin struct {
	mutex    sync.RWMutex
	Counter  int
	Backends []*backend.Backend
}

func (rr *RoundRobin) AddBackend(b *backend.Backend) {
	rr.Backends = append(rr.Backends, b)
}

func (rr *RoundRobin) GetBackends() []*backend.Backend {
	return rr.Backends
}

func (rr *RoundRobin) nextServer() (*backend.Backend, error) {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	attempts := 1
	for {
		b := rr.Backends[rr.Counter]
		rr.Counter = (rr.Counter + 1) % len(rr.Backends)
		if b.Alive {
			return b, nil
		}
		fmt.Println()
		if attempts >= len(rr.Backends) {
			fmt.Println("AHHHHH, none is alive!!!!")
			return nil, ErrNoServer
		}
		attempts++
	}
}

func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := rr.nextServer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	b.Proxy.ServeHTTP(w, r)
}

func New() *RoundRobin {
	return &RoundRobin{}
}

func NewWithBackends(backends []*backend.Backend) *RoundRobin {
	return &RoundRobin{
		Backends: backends,
	}
}
