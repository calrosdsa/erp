package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type ReceiptEventData struct {
	Body             dto.ReceiptBody
	Receipt          model.Receipt
	ReceiptPartyType string
	Tx               *query.QueryTx
}

type StatusReceiptEventData struct {
	Receipt          model.Receipt
	ReceiptPartyType string
	CompanyDefault   model.CompanyDefault
	Company          model.Company
	StockDefault     model.StockDefault
	LineItemsData    dto.LineItemsData
	Tx               *query.QueryTx
}
