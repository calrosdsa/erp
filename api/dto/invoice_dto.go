package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	RequestInvoices struct {
		PartyType string `path:"party" required:"true"`
		PaginationParams
		OptionalQueryParams

		DueDate      string `query:"due_date" required:"false"`
		PostingDate  string `query:"posting_date" required:"false"`
		PartyID      string `query:"party_id" required:"false"`
		Code         string `query:"code" required:"false"`
		ProjectID    string `query:"project_id" required:"false"`
		CostCenterId string `query:"cost_center_id" required:"false"`

		OrderID string `query:"order_id" required:"false"`
	}
	CreateInvoiceRequest struct {
		Body InvoiceBody
	}
	EditInvoiceRequest struct {
		Body InvoiceBody
	}

	InvoiceBody struct {
		Invoice             InvoiceData         `json:"invoice"`
		CreateItemLines     CreateItemLines     `json:"items"`
		CreateTaxAndCharges CreateTaxAndChanges `json:"tax_and_charges"`
	}

	InvoiceData struct {
		ID               int64   `json:"id" required:"false"`
		InvoicePartyType string  `json:"invoice_party_type" required:"true"`
		TotalAmount      float64 `json:"total_amount" required:"true"`
		RecordNo         string  `json:"record_no" required:"false"`

		Fields InvoiceFields `json:"fields"`
	}
	InvoiceFields struct {
		DueDate     *time.Time `json:"due_date" required:"false"`
		PostingDate time.Time  `json:"posting_date" required:"true"`
		PostingTime string     `json:"posting_time"`
		Tz          string     `json:"tz"`
		Currency    string     `json:"currency" required:"true"`
		PartyID     int64      `json:"party_id" required:"true"`

		UpdateStock bool `json:"update_stock" required:"false"`

		DocReferenceID *int64 `json:"doc_reference_id" required:"false"`

		ProjectID    *int64 `json:"project_id" required:"false"`
		CostCenterID *int64 `json:"cost_center_id" required:"false"`
		PriceListID  *int64 `json:"price_list_id" required:"false"`
		WarehouseID  *int64 `json:"warehouse_id" required:"false"`
	}

	InvoiceDetailDto struct {
		Invoice        InvoiceDto             `json:"invoice"`
		Totals         TotalsDto              `json:"totals"`
		AcctDimensions AccountingDimensionDto `json:"acct_dimensions"`
	}

	TotalsDto struct {
		Paid  int32 `json:"paid"`
		Total int32 `json:"total"`
	}

	InvoiceDto struct {
		ID             int64      `json:"id"`
		Code           string     `json:"code"`
		DueDate        *time.Time `json:"due_date"`
		Date           time.Time  `json:"date"`
		CreatedAt      time.Time  `json:"created_at"`
		Currency       string     `json:"currency"`
		Status         string     `json:"status"`
		PostingDate    time.Time  `json:"posting_date"`
		PostingTime    string     `json:"posting_time"`
		Tz             string     `json:"tz"`
		UpdateStock    bool       `json:"update_stock"`
		RecordNo       string     `json:"record_no"`
		PaidAmount     int32      `json:"paid_amount"`
		TotalAmount    int32      `json:"total_amount"`
		DocReferenceID *int64     `json:"doc_reference_id"`

		PriceListField
		WarehouseOptionalField
		AccountingDimensionDto
		TotalsDto
		//Can be supplier or customer
		PartyID   int64  `json:"party_id"`
		PartyType string `json:"party_type"`
		PartyName string `json:"party_name"`
		PartyUUID string `json:"party_uuid"`
	}
)

func InvoiceDtoFromModel(m *model.Invoice) InvoiceDto {
	r := InvoiceDto{}
	r.ID = m.ID
	r.Code = m.Code
	r.DueDate = m.DueDate
	r.CreatedAt = m.CreatedAt
	r.Status = m.Status
	r.Currency = m.Currency
	return r
}
