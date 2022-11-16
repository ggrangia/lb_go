package lb_go

import (
	"fmt"
	"net/http"

	"github.com/ggrangia/lb_go/pkg/backend"
	"github.com/ggrangia/lb_go/pkg/selection"
)

type Lb struct {
	Backends []backend.Backend
	Selector selection.Selector
}

func (lb *Lb) Serve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called the proxy")
	b := lb.Selector.Select(lb.Backends)

	b.Proxy.ServeHTTP(w, r)
}

func (lb *Lb) Start() {

}
