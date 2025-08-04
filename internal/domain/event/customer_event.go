package event

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type CustomerEventData struct {
	Tx       *query.QueryTx
	Customer model.Customer
	Data     dto.CustomerData
	Req      common.RequestContext
}
