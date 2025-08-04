package event

import (
	"erp/api/dto"
	// "erp/gen/db/model"
	"erp/gen/db/query"
)

type StatusPaymentTermsTemplateEventData struct {
	PaymentTermsTemplateID int64
	Tx                 *query.QueryTx
}

type PaymentTermsTemplateEventData struct {
	PaymentTermsTemplateID int64
	Tx        *query.QueryTx
	Body      dto.PaymentTermsTemplateData
}
