package supplier_rest

type SupplierPaths struct {
	Base string
	Detail string 
	UpdateStatus string
}

func NewSupplierPaths(base string)SupplierPaths {
	return SupplierPaths{
		Base: base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
	}
}

type PurchasePaths struct {
	Base string 
	Detail string
}

func NewPurchasePaths(base string) PurchasePaths{
	return PurchasePaths{
		Base:base,
		Detail : base + "/detail/{id}",
	}
}