package lb_go

import "github.com/ggrangia/lb_go/pkg/backend"

type Selector interface {
	Select([]backend.Backend) backend.Backend
}

type Lb struct {
	Backends []backend.Backend
	Selector Selector
}
