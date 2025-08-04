package rest_price_list

type PriceListPaths struct {
	Base   string
	Detail string
}

func NewPriceListPaths(base string) PriceListPaths {
	return PriceListPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
	}
}
