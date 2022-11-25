package selection

import (
	"net/http"

	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
)

type Selector interface {
	Select() *backend.Backend
	http.Handler
}
