package order_rest

type OrderPaths struct {
	Base   string
	Type   string
	Detail string
	UpdateStatus string
	Document string 
}

func NewOrderPaths(base string) OrderPaths {
	return OrderPaths{
		Base:   base,
		Type:   base + "/{party}",
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
		Document: base + "/export/document",
	}
}
