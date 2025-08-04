package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type StatusQuotationEventData struct {
	Quotation          model.Quotation
	QuotationPartyType string
	LineItemsData      dto.LineItemsData
	Tx                 *query.QueryTx
}

type QuotationEventData struct {
	Quotation *model.Quotation
	Tx        *query.QueryTx
	Body      dto.QuotationBody
}
