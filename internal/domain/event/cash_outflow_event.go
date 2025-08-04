package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type CashOutflowEventData struct {
	Data        dto.CashOutflowData
	CashOutflow model.CashOutflow
	Tx          *query.QueryTx
}

type StatusCashOutflowEventData struct {
	CashOutflow     model.CashOutflow
	Tx              *query.QueryTx
	CompanyDefaults model.CompanyDefault
	TaxLinesData     dto.TaxLinesData
}
