package dto

import "erp/gen/db/model"

type (
	PaymentTermsRequest struct {
		DefaultListParams
		Name      string `query:"name"`
		CreatedAt string `query:"created_at"`
	}

	PaymentTermsDataRequest struct {
		Body PaymentTermsData
	}
	PaymentTermsData struct {
		ID     int64              `json:"id" required:"false"`
		Fields PaymentTermsFields `json:"filds" required:"true"`
	}
	PaymentTermsFields struct {
		Name                   string  `json:"name"`
		InvoicePortion         int32   `json:"invoice_portion"`
		CreditDays             *int32  `json:"credit_days"`
		DueDateBaseOn          string  `json:"due_date_base_on"`
		Description            *string `json:"description"`
		DiscountType           *string `json:"discount_type"`
		Discount               *int64  `json:"discount"`
		DiscountValidityBaseOn *string `json:"discount_validity_base_on"`
		DiscountValidity       *int32  `json:"discount_validity"`
	}
	PaymentTermsDto struct {
		ID                     int64   `json:"id"`
		UUID                   string  `json:"uuid"`
		Status                 string  `json:"status"`
		Name                   string  `json:"name"`
		InvoicePortion         int32   `json:"invoice_portion"`
		CreditDays             *int32  `json:"credit_days"`
		DueDateBaseOn          string  `json:"due_date_base_on"`
		Description            *string `json:"description"`
		DiscountType           *string `json:"discount_type"`
		Discount               *int64  `json:"discount"`
		DiscountValidityBaseOn *string `json:"discount_validity_base_on"`
		DiscountValidity       *int32  `json:"discount_validity"`
	}

	PaymentTermsLineDto struct {
		ID             int32   `json:"id"`

		PaymentTermsID int64   `json:"payment_term_id"`
		PaymentTerm    string   `json:"payment_term"`
		
		DocumentID     int64   `json:"document_id"`
		InvoicePortion int32   `json:"invoice_portion"`
		CreditDays     *int32  `json:"credit_days"`
		DueDateBaseOn  string  `json:"due_date_base_on"`
		Description    *string `json:"description"`
	}

	PaymentTermsLineData struct {
		PaymentTermsID int64   `json:"payment_terms_id"`
		InvoicePortion int32   `json:"invoice_portion"`
		CreditDays     *int32  `json:"credit_days"`
		DueDateBaseOn  string  `json:"due_date_base_on"`
		Description    *string `json:"description"`
	}
)

func PaymentTermsFromModel(m model.PaymentTerm) PaymentTermsDto {
	return PaymentTermsDto{
		ID:   m.ID,
		UUID: m.UUID,
		Name: m.Name,
	}
}
