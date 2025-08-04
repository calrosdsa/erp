package selling

type SalesOrderPaths struct {
	Base         string
	DetailOrder  string
	ClientOrders string
}

func NewSalesOrderPaths(base string) SalesOrderPaths {
	return SalesOrderPaths{
		Base:         base,
		ClientOrders: base + "/client",
		DetailOrder:  base + "/{code}",
	}
}

type TaxPaths struct {
	Base   string
	Detail string
}

func NewTaxPaths(base string) TaxPaths {
	return TaxPaths{
		Base:   base,
		Detail: base + "/{id}",
	}
}

type CustomerPaths struct {
	Base   string
	CustomerTypes string
	Detail string
}

func NewCustomerPaths(base string) CustomerPaths {
	return CustomerPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
		CustomerTypes: base + "/customer-types",
	}
}
