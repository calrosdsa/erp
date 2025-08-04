package repository

import (
	"erp/api/common"
	"erp/gen/db/model"
)

type StockService interface {
	GetStockDefault(req *common.RequestContext) (model.StockDefault, error)
	// GetLineItems(tx *query.QueryTx, req *common.RequestContext, id int64, opts ...OptionStock) (dto.LineItemsData, error)
}
