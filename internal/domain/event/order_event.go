package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type OrderEventData struct {
	Order model.Order
	Tx    *query.QueryTx
	Body  dto.OrderBody
	
}
