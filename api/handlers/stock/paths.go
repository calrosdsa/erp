package stock

type ItemGroupPath struct {
	Base string
	Detail string
}

func NewItemGroupPath(base string) ItemGroupPath {
	return ItemGroupPath{
		Base: base,
		Detail: base + "/{id}",
	}
}

type ItemPath struct {
	Base          string
	Detail string
	ItemPriceList string
	ItemPrice     string
}

func NewItemPath(base string) ItemPath {
	return ItemPath{
		Base:          base,
		Detail: base + "/{id}",
	}
}


type ItemPricePaths struct {
	Base string
	Item string
	Detail string 
	Order string
}

func NewItemPricePaths(base string) ItemPricePaths {
	return ItemPricePaths{
		Base: base,
		Item: base + "/{itemCode}",
		Detail: base + "/detail/{id}",
		Order: base + "/order",
	}
}



type PriceListPath struct {
	Base string
	Detail string
}

func NewPriceListPaths(base string) PriceListPath {
	return PriceListPath{
		Base: base,
		Detail: base + "/{id}",
	}
}



type ItemAttributePaths struct {
	Base string
	Detail string
	ItemAttributeValue string
}

func NewItemAttributePaths(base string) ItemAttributePaths {
	return ItemAttributePaths{
		Base: base,
		Detail: base + "/{id}",
		ItemAttributeValue: base + "/item-attribute-value",
	}
}


type ItemVariantPaths struct {
	Base string
}
func NewItemVariantPaths(base string)ItemVariantPaths {
	return ItemVariantPaths{
		Base: base,
	}
}

type ItemStockPaths struct {
	Base string
	Item string 
	Warehouse string 
}

func NewItemStockPaths(base string)ItemStockPaths {
	return ItemStockPaths{
		Base: base,
		Item: base + "/item",
		Warehouse: base + "/warehouse",
	}
}