package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type ChargesTemplateEventData struct {
	Tx                 *query.QueryTx
	ChargesTemplate model.ChargesTemplate
	ChargeTemplateData dto.ChargesTemplateBody
}