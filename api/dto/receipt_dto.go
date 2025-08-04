package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	RequestReceipts struct {
		PartyType string `path:"party" required:"true"`
		PaginationParams
		OptionalQueryParams
		PostingDate string `query:"posting_date" required:"false"`
		PartyID     string `query:"party_id" required:"false"`
	}

	CreateReceiptRequest struct {
		Body ReceiptBody
	}
	EditReceiptRequest struct {
		Body ReceiptBody
	}

	ReceiptBody struct {
		Receipt             ReceiptData         `json:"receipt" required:"true"`
		CreateItemLines     CreateItemLines     `json:"items"`
		CreateTaxAndCharges CreateTaxAndChanges `json:"tax_and_charges"`
	}

	ReceiptData struct {
		ID               int64         `json:"id" requried:"false"`
		ReceiptPartyType string        `json:"receipt_party_type" required:"true"`
		Fields           ReceiptFields `json:"fields"`
	}

	ReceiptFields struct {
		PostingDate time.Time  `json:"posting_date" required:"true"`
		PostingTime string     `json:"posting_time"`
		Tz          string     `json:"tz"`
		Currency    string     `json:"currency" required:"true"`
		PartyID     int64      `json:"party_id" required:"true"`
		WarehouseID  int64 `json:"warehouse_id"`

		DocReferenceID *int64 `json:"doc_reference_id"`

		ProjectID    *int64 `json:"project_id" required:"false"`
		CostCenterID *int64 `json:"cost_center_id" required:"false"`
		PriceListID  *int64 `json:"price_list_id" required:"false"`
	}

	CreateItemLines struct {
		Lines       []LineItemData `json:"lines" required:"true"`
	}

	ReceiptDto struct {
		ID          int64     `json:"id"`
		Code        string    `json:"code"`
		Currency    string    `json:"currency"`
		PostingDate time.Time `json:"posting_date"`
		PostingTime string    `json:"posting_time"`
		Tz          string    `json:"tz"`
		CreatedAt   time.Time `json:"created_at"`
		Status      string    `json:"status"`

		PartyID   int64  `json:"party_id"`
		PartyName string `json:"party_name"`
		PartyUUID string `json:"party_uuid"`
		PartyType string `json:"party_type"`

		DocReferenceID *int64 `json:"doc_reference_id"`

		PriceListField
		WarehouseField

		AccountingDimensionDto
	}
	ReceiptDetailDto struct {
		Receipt   ReceiptDto    `json:"receipt"`
		ItemLines []ItemLineDto `json:"item_lines"`
	}
)

func ReceiptDtoFromModel(m *model.Receipt) ReceiptDto {
	r := ReceiptDto{}
	r.Code = m.Code
	r.Currency = m.Currency
	r.CreatedAt = m.CreatedAt
	r.PostingDate = m.PostingDate
	return r
}
