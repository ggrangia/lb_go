package randomselection

import (
	"errors"
	"math/rand"
	"net/http"
	"sync"

	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
)

var ErrNoServer = errors.New("no available servers")

type RandomSelection struct {
	generator rand.Rand
	Backends  []*backend.Backend
	mutex     sync.RWMutex
}

func (rs *RandomSelection) nextServer() (*backend.Backend, error) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	l := len(rs.Backends)
	// array length cannot be less than zero
	if l <= 0 {
		return nil, ErrNoServer
	}
	return rs.Backends[rs.generator.Intn(l)], nil

}

func (rs *RandomSelection) GetBackends() []*backend.Backend {
	return rs.Backends
}

func New(seed int64) *RandomSelection {
	source := rand.NewSource(seed)
	generator := rand.New(source)
	return &RandomSelection{
		generator: *generator,
	}
}

func NewWithBackends(seed int64, backends []*backend.Backend) *RandomSelection {
	rs := New(seed)
	rs.Backends = backends

	return rs
}

func (rs *RandomSelection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := rs.nextServer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b.Proxy.ServeHTTP(w, r)
}

func (rs *RandomSelection) AddBackend(b *backend.Backend) {
	rs.Backends = append(rs.Backends, b)
}
