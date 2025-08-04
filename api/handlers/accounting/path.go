package accounting


type TaxPaths struct {
	Base string
	Detail string
}

func NewTaxPaths(base string) TaxPaths {
	return TaxPaths{
		Base: base,
		Detail: base + "/{id}",
	}
}

// type ItemPath struct {
// 	Base          string
// 	ItemPriceList string
// 	ItemPrice     string
// }

// func NewItemPath(base string) ItemPath {
// 	return ItemPath{
// 		Base:          base,
// 		ItemPriceList: base + "/item-price-list",
// 		ItemPrice:     base + "/item-price",
// 	}
// }

// type PriceListPath struct {
// 	Base string
// }

// func NewPriceListPaths(base string) PriceListPath {
// 	return PriceListPath{
// 		Base: base,
// 	}
// }
