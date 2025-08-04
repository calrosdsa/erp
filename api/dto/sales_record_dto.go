package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	ExportDataRequest struct {
		Body struct {
			Data map[string]string `json:"data"`
		}
	}
	SalesRecordsRequest struct {
		PaginationParams
		OptionalQueryParams
		InvoiceDate string `query:"invoice_date"`
		InvoiceID   string `query:"invoice_id"`
		CustomerID string `query:"customer_id"`
	}
	CreateSalesRecordRequest struct {
		Body SalesRecordData

	}
	EditSalesRecordRequest struct {
		Body SalesRecordData
	}

	SalesRecordData struct {
		ID                                   int64     `json:"id" required:"false"`
		InvoiceDate                          time.Time `json:"invoice_date"`
		InvoiceNo                            string    `json:"invoice_no"`
		AuthorizationCode                    string    `json:"authorization_code"`
		CustomerNitCi                        string    `json:"customer_nit_ci"`
		Supplement                           string    `json:"supplement"`
		NameOrBusinessName                   string    `json:"name_or_business_name"`
		TotalSaleAmount                      float64   `json:"total_sale_amount"`
		IceAmount                            float64   `json:"ice_amount"`
		IehdAmount                           float64   `json:"iehd_amount"`
		IpjAmount                            float64   `json:"ipj_amount"`
		TaxRates                             float64   `json:"tax_rates"`
		OtherNotSubjectToVat                 float64   `json:"other_not_subject_to_vat"`
		ExportsAndExemptOperations           float64   `json:"exports_and_exempt_operations"`
		ZeroRateTaxableSales                 float64   `json:"zero_rate_taxable_sales"`
		Subtotal                             float64   `json:"subtotal"`
		DiscountsBonusAndRebatesSubjectToVat float64   `json:"discounts_bonus_and_rebates_subject_to_vat"`
		GiftCardAmount                       float64   `json:"gift_card_amount"`
		BaseAmountForTaxDebit                float64   `json:"base_amount_for_tax_debit"`
		TaxDebit                             float64   `json:"tax_debit"`
		State                                string    `json:"state"`
		ControlCode                          string    `json:"control_code"`
		SaleType                             string    `json:"sale_type"`
		WithTaxCreditRight                   bool      `json:"with_tax_credit_right"`
		ConsolidationStatus                  string    `json:"consolidation_status"`
		CustomerID                           int64     `json:"customer_id"`
		InvoiceID                            int64     `json:"invoice_id"`
	}
	SalesRecordDto struct {
		ID                                   int64     `json:"id"`
		UUID                                 string    `json:"uuid"`
		InvoiceDate                          time.Time `json:"invoice_date"`
		InvoiceNo                            string    `json:"invoice_no"`
		AuthorizationCode                    string    `json:"authorization_code"`
		CustomerNitCi                        string    `json:"customer_nit_ci"`
		Supplement                           string    `json:"supplement"`
		NameOrBusinessName                   string    `json:"name_or_business_name"`
		TotalSaleAmount                      int32     `json:"total_sale_amount"`
		IceAmount                            int32     `json:"ice_amount"`
		IehdAmount                           int32     `json:"iehd_amount"`
		IpjAmount                            int32     `json:"ipj_amount"`
		TaxRates                             int32     `json:"tax_rates"`
		OtherNotSubjectToVat                 int32     `json:"other_not_subject_to_vat"`
		ExportsAndExemptOperations           int32     `json:"exports_and_exempt_operations"`
		ZeroRateTaxableSales                 int32     `json:"zero_rate_taxable_sales"`
		Subtotal                             int32     `json:"subtotal"`
		DiscountsBonusAndRebatesSubjectToVat int32     `json:"discounts_bonus_and_rebates_subject_to_vat"`
		GiftCardAmount                       int32     `json:"gift_card_amount"`
		BaseAmountForTaxDebit                int32     `json:"base_amount_for_tax_debit"`
		TaxDebit                             int32     `json:"tax_debit"`
		State                                string    `json:"state"`
		ControlCode                          string    `json:"control_code"`
		SaleType                             string    `json:"sale_type"`
		WithTaxCreditRight                   bool      `json:"with_tax_credit_right"`
		ConsolidationStatus                  string    `json:"consolidation_status"`
		Customer                             string    `json:"customer"`
		CustomerID                           int64     `json:"customer_id"`
		CustomerUUID                         string    `json:"customer_uuid"`
		Status                               string    `json:"status"`
		InvoiceCode                          string    `json:"invoice_code"`
		InvoiceID                            int64     `json:"invoice_id"`
	}
)

func SalesRecordDtoFromModel(m *model.SalesRecord) SalesRecordDto {
	return SalesRecordDto{
		ID:        m.ID,
		UUID:      m.UUID,
		InvoiceNo: m.InvoiceNo,
	}
}
