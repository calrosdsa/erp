package event

import (
	"erp/api/common"
	"erp/api/dto"

	"erp/gen/db/query"
)

type ItemCreatedEventData struct {
	Item *dto.ItemDto
	Req  common.RequestContext
	Tx   *query.QueryTx
	Body dto.ItemData
}
