package event

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type DealEventData struct {
	Data dto.DealData
	Req  common.RequestContext
	Deal model.Deal
	Tx   *query.QueryTx
}

type ChangeStageEventData struct {
	ProfileID       int64
	StageTransition dto.EntityTransitionData
	Tx              *query.QueryTx
}
