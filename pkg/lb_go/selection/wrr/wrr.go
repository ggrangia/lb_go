package wrr

import (
	"container/heap"
	"errors"
	"net/http"
	"sync"

	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
)

var ErrNoServer = errors.New("no available servers")

// Wrr based on Earliest Deadline First (floating point)
// The deadline is computed as current + 1 / weight
// current simulates the flow of time

type weightedBackend struct {
	backend.Backend
	weight   int
	deadline float64
}

type Wrr struct {
	mutex    sync.RWMutex
	Backends []*weightedBackend
	current  float64
}

func New() *Wrr {
	return &Wrr{}
}

func (w *Wrr) AddWeightedBackend(b *backend.Backend, i int) {
	if i <= 0 {
		// meaningless
		return
	}
	w.mutex.Lock()
	d := w.current + 1/float64(i)
	wb := &weightedBackend{*b, i, d}
	heap.Push(w, wb)
	w.mutex.Unlock()
}

func (w *Wrr) Len() int {
	return len(w.Backends)
}

func (w *Wrr) Less(i, j int) bool {
	return w.Backends[i].deadline < w.Backends[j].deadline
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
	attempts := 1
	for {
		wb := heap.Pop(w).(*weightedBackend)
		// update current time
		w.current = wb.deadline
		// Update the backend (using new current) and put it back in the heap
		wb.deadline = w.current + 1/float64(wb.weight)
		heap.Push(w, wb)

		if wb.Alive {
			return wb, nil
		}
		if attempts >= len(w.Backends) {
			return nil, ErrNoServer
		}
		attempts++
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
