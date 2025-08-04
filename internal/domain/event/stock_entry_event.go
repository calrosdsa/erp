package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type StockEntryEventData struct {
	StockEntry     model.StockEntry
	Tx             *query.QueryTx
	StockEntryBody dto.StockEntryBody
}

type StatusStockEntryEventData struct {
	StockEntry      model.StockEntry
	Tx              *query.QueryTx
	CompanyID       int64
	LineItemsData   dto.LineItemsData
	CompanyDefaults model.CompanyDefault
}
