package wrr

import (
	"container/heap"
	"errors"
	"net/http"
	"sync"

	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
)

var ErrNoServer = errors.New("no available servers")

type weightedBackend struct {
	backend.Backend
	weight int
}

type Wrr struct {
	mutex    sync.RWMutex
	Backends []*weightedBackend
}

func New() *Wrr {
	return &Wrr{}
}

func (w *Wrr) AddWeightedBackend(b *backend.Backend, i int) {
	w.mutex.Lock()
	wb := &weightedBackend{*b, i}
	heap.Push(w, wb)
	w.mutex.Unlock()
}

func (w *Wrr) Len() int {
	return len(w.Backends)
}

func (w *Wrr) Less(i, j int) bool {
	return w.Backends[i].weight < w.Backends[j].weight
}

func (w *Wrr) Swap(i, j int) {
	w.Backends[i], w.Backends[j] = w.Backends[j], w.Backends[i]
}

func (w *Wrr) Push(x any) {
	item, ok := x.(*weightedBackend)
	if !ok {
		return
	}
	w.Backends = append(w.Backends, item)
}

func (w *Wrr) Pop() interface{} {
	l := len(w.Backends)
	b := w.Backends[l-1]
	w.Backends[l-1] = nil
	w.Backends = w.Backends[0 : l-1]
	return b
}

func (w *Wrr) GetBackends() []*backend.Backend {
	backends := make([]*backend.Backend, len(w.Backends))
	for i, b := range w.Backends {
		backends[i] = &b.Backend
	}
	return backends
}

func (w *Wrr) nextServer() (*weightedBackend, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(w.Backends) == 0 {
		return nil, ErrNoServer
	}
	for {
		// TODO: Implement
	}
}

func (w *Wrr) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	b, err := w.nextServer()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	b.Proxy.ServeHTTP(rw, r)
}
