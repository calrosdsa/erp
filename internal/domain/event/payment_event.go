package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type PaymentEventData struct {
	Body    dto.PaymentBody
	Payment model.Payment
	Tx      *query.QueryTx
}

type StatusPaymentEventData struct {
	Payment         model.Payment
	Tx              *query.QueryTx
	CompanyDefaults model.CompanyDefault
	TaxLinesData    dto.TaxLinesData
	References      []*model.PaymentReference
}
