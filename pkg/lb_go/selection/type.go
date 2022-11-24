package selection

import "github.com/ggrangia/lb_go/pkg/backend"

type Selector interface {
	Select() backend.Backend
}
