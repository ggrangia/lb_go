package lb_go

type Lb struct {
	Backends []Backend
	Selector Selector
}

type Selector interface {
	Select([]Backend) Backend
}
