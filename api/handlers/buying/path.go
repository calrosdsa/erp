package buying

type SupplierPaths struct {
	Base string
	Detail string 
}

func NewSupplierPaths(base string)SupplierPaths {
	return SupplierPaths{
		Base: base,
		Detail: base + "/detail/{id}",
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