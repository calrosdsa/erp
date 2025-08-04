package customer_rest

type CustomerPaths struct {
	Base          string
	CustomerTypes string
	Detail        string
	UpdateStatus  string
}

func NewCustomerPaths(base string) CustomerPaths {
	return CustomerPaths{
		Base:          base,
		Detail:        base + "/detail/{id}",
		CustomerTypes: base + "/customer-types",
		UpdateStatus:  base + "/update-status",
	}
}
