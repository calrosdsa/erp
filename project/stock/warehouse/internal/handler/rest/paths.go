package warehouse_rest

type WareHousePaths struct {
	Base     string
	TreeView string
	Detail   string
}

func NewWareHousePaths(base string) WareHousePaths {
	return WareHousePaths{
		Base:     base,
		Detail:   base + "/detail/{id}",
		TreeView: base + "/view/tree",
	}
}
