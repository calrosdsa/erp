package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type StatusInvoiceEventData struct {
	Invoice          model.Invoice
	InvoicePartyType string
	CompanyDefaults  model.CompanyDefault
	LineItemsData    dto.LineItemsData
	TaxLinesData     dto.TaxLinesData
	Tx               *query.QueryTx
}

type InvoiceEventData struct {
	Invoice          model.Invoice
	InvoicePartyType string
	Tx               *query.QueryTx
	Body             dto.InvoiceBody
}
