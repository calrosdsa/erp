package itemline_rest

type ItemLinePaths struct {
	Base        string
	ProductList string
}

func NewItemLinePaths(base string) ItemLinePaths {
	return ItemLinePaths{
		Base:        base,
		ProductList: base + "/products",
	}
}
