package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	PurchaseRecordsRequest struct {
		PaginationParams
		OptionalQueryParams
		InvoiceDate string `query:"invoice_date"`
		InvoiceID   string `query:"invoice_id"`
		SupplierID  string `query:"supplier_id"`
	}
	CreatePurchaseRecordRequest struct {
		Body PurchaseRecordData
	}
	EditPurchaseRecordRequest struct {
		Body PurchaseRecordData
	}

	PurchaseRecordData struct {
		ID     int64                `json:"id"`
		Fields PurchaseRecordFields `json:"fields" required:"true"`
	}

	PurchaseRecordFields struct {
		SupplierNit                       string    `json:"supplier_nit"`
		SupplierBusinessName              string    `json:"supplier_business_name"`
		AuthorizationCode                 string    `json:"authorization_code"`
		InvoiceNo                         string    `json:"invoice_no"`
		DuiDimNo                          string    `json:"dui_dim_no"`
		InvoiceDuiDimDate                 time.Time `json:"invoice_dui_dim_date"`
		TotalPurchaseAmount               int64     `json:"total_purchase_amount"`
		IceAmount                         int64     `json:"ice_amount"`
		IehdAmount                        int64     `json:"iehd_amount"`
		IpjAmount                         int64     `json:"ipj_amount"`
		TaxRates                          int32     `json:"tax_rates"`
		OtherNotSubjectToTaxCredit        int64     `json:"other_not_subject_to_tax_credit"`
		ExemptAmounts                     int64     `json:"exempt_amounts"`
		ZeroRateTaxablePurchasesAmount    int64     `json:"zero_rate_taxable_purchases_amount"`
		Subtotal                          int64     `json:"subtotal"`
		DiscountsBonusRebatesSubjectToVat int64     `json:"discounts_bonus_rebates_subject_to_vat"`
		GiftCardAmount                    int64     `json:"gift_card_amount"`
		CfBaseAmount                      int64     `json:"cf_base_amount"`
		TaxCredit                         int64     `json:"tax_credit"`
		PurchaseType                      string    `json:"purchase_type"`
		ControlCode                       string    `json:"control_code"`
		WithTaxCreditRight                bool      `json:"with_tax_credit_right"`
		ConsolidationStatus               string    `json:"consolidation_status"`
		SupplierID                        int64     `json:"supplier_id"`
		InvoiceID                         *int64    `json:"invoice_id" required:"false"`
	}
	PurchaseRecordDto struct {
		ID                                int64     `json:"id"`
		UUID                              string    `json:"uuid"`
		SupplierNit                       string    `json:"supplier_nit"`
		SupplierBusinessName              string    `json:"supplier_business_name"`
		AuthorizationCode                 string    `json:"authorization_code"`
		InvoiceNo                         string    `json:"invoice_no"`
		DuiDimNo                          string    `json:"dui_dim_no"`
		InvoiceDuiDimDate                 time.Time `json:"invoice_dui_dim_date"`
		TotalPurchaseAmount               int64     `json:"total_purchase_amount"`
		IceAmount                         int64     `json:"ice_amount"`
		IehdAmount                        int64     `json:"iehd_amount"`
		IpjAmount                         int64     `json:"ipj_amount"`
		TaxRates                          int32     `json:"tax_rates"`
		OtherNotSubjectToTaxCredit        int64     `json:"other_not_subject_to_tax_credit"`
		ExemptAmounts                     int64     `json:"exempt_amounts"`
		ZeroRateTaxablePurchasesAmount    int64     `json:"zero_rate_taxable_purchases_amount"`
		Subtotal                          int64     `json:"subtotal"`
		DiscountsBonusRebatesSubjectToVat int64     `json:"discounts_bonus_rebates_subject_to_vat"`
		GiftCardAmount                    int64     `json:"gift_card_amount"`
		CfBaseAmount                      int64     `json:"cf_base_amount"`
		TaxCredit                         int64     `json:"tax_credit"`
		PurchaseType                      string    `json:"purchase_type"`
		ControlCode                       string    `json:"control_code"`
		WithTaxCreditRight                bool      `json:"with_tax_credit_right"`
		ConsolidationStatus               string    `json:"consolidation_status"`
		SupplierID                        int64     `json:"supplier_id"`
		Supplier                          string    `json:"supplier"`
		SupplierUUID                      string    `json:"supplier_uuid"`
		Status                            string    `json:"status"`
		CreatedAt                         time.Time `json:"created_at"`

		InvoiceCode *string `json:"invoice_code"`
		InvoiceID   *int64  `json:"invoice_id"`
	}
)

func PurchaseRecordDtoFromModel(m *model.PurchaseRecord) PurchaseRecordDto {
	return PurchaseRecordDto{
		ID:        m.ID,
		UUID:      m.UUID,
		InvoiceNo: m.InvoiceNo,
	}
}
