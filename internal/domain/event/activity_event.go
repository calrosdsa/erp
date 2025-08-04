package event

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
)

type ActivityEventData struct {
	Tx *query.QueryTx
	Data dto.ActivityData
	ReqCtx common.RequestContext
}