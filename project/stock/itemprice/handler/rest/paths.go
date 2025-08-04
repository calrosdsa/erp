package itemprice_rest


type ItemPricePaths struct {
	Base string
	Item string
	Detail string 
	Order string
	AssociatedActions string
}

func NewItemPricePaths(base string) ItemPricePaths {
	return ItemPricePaths{
		Base: base,
		Item: base + "/{itemCode}",
		Detail: base + "/detail/{id}",
		Order: base + "/order",
		AssociatedActions: base + "/associated-actions",
	}
}
