package dto

import "erp/gen/db/model"

type (
	PaymentTermsTemplateRequest struct {
		DefaultListParams
		Name      string `query:"name"`
		CreatedAt string `query:"created_at"`
	}

	PaymentTermsTemplateDataRequest struct {
		Body PaymentTermsTemplateData
	}
	PaymentTermsTemplateData struct {
		ID     int64              `json:"id" required:"false"`
		Lines []PaymentTermsLineData `json:"lines" required:"true"`
		Fields PaymentTermsTemplateFields `json:"filds" required:"true"`
	}
	PaymentTermsTemplateFields struct {
		Name                   string  `json:"name"`
	}
	PaymentTermsTemplateDto struct {
		ID                     int64   `json:"id"`
		UUID                   string  `json:"uuid"`
		Status                 string  `json:"status"`
		Name                   string  `json:"name"`
	}
)

func PaymentTermsTemplateFromModel(m model.PaymentTermsTemplate) PaymentTermsTemplateDto {
	return PaymentTermsTemplateDto{
		ID:           m.ID,
		UUID:         m.UUID,
		Name: m.Name,
	}
}
