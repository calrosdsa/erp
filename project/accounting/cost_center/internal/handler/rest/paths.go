package rest_cost_center

type CostCenterPaths struct {
	Base   string
	Detail string
}

func NewCostCenterPaths(base string) CostCenterPaths {
	return CostCenterPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
	}
}
