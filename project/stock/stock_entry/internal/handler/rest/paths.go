package rest_stock_entry

type StockEntryPaths struct {
	Base         string
	Detail       string
	UpdateStatus string
}

func NewStockEntryPaths(base string) StockEntryPaths {
	return StockEntryPaths{
		Base:         base,
		Detail:       base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
	}
}
