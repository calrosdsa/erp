package event

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type TaskEventData struct {
	Data dto.TaskData
	Req  common.RequestContext
	Task model.Task
	Tx   *query.QueryTx
}


