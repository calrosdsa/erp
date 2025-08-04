package event

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
)

type SupplierEventData struct {
	Tx       *query.QueryTx
	Supplier dto.SupplierDto
	Data     dto.SupplierData
	Req      common.RequestContext
}

// type CustomerEventData struct {
// 	Tx       *query.QueryTx
// 	Customer model.Customer
// 	Data     dto.CustomerData
// 	Req      common.RequestContext
// }
