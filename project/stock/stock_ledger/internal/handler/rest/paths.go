package stock_ledger_rest

type StockLedgerPaths struct {
	StockLedgerReport  string
	StockBalanceReport  string
}

func NewStockLedgerPaths(base string) StockLedgerPaths {
	return StockLedgerPaths{
		StockLedgerReport:  base + "/report",
		StockBalanceReport:base + "/balance",
	}
}
