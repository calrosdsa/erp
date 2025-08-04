package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	RequestQuotations struct {
		PartyType string `path:"party" required:"true"`
		PaginationParams
		OptionalQueryParams
		ValidTill    string `query:"valid_till" required:"false"`
		PostingDate  string `query:"posting_date" required:"false"`
		PartyID      string `query:"party_id" required:"false"`
		ID           string `query:"id" required:"false"`
		ProjectID    string `query:"project_id" required:"false"`
		CostCenterID string `query:"cost_center_id" required:"false"`
		PricingID    string `query:"pricing_id" required:"false"`
	}

	CreateQuotationRequest struct {
		Body QuotationBody
	}

	EditQuotationRequest struct {
		Body QuotationBody
	}

	QuotationBody struct {
		Quotation           QuotationData       `json:"quotation" required:"true"`
		CreateItemLines     CreateItemLines     `json:"items" required:"true"`
		CreateTaxAndCharges CreateTaxAndChanges `json:"tax_and_charges"`
	}

	QuotationData struct {
		ID         int64           `json:"id" required:"false"`
		QuotationPartyType string    `json:"quotation_party_type" required:"true"`
		Fields     QuotationFields `json:"fields"`
		References []*int64         `json:"references" required:"false"`
	}

	QuotationFields struct {
		Currency           string    `json:"currency"`
		PostingDate        time.Time `json:"posting_date"`
		PostingTime        string    `json:"posting_time"`
		Tz                 string    `json:"tz"`
		PartyID            int64     `json:"party_id" required:"true"`
		ValidTill          time.Time `json:"valid_till" required:"true"`
		
		ProjectID    *int64 `json:"project_id" required:"false"`
		CostCenterID *int64 `json:"cost_center_id" required:"false"`
		PriceListID  *int64     `json:"price_list_id" required:"false"`
	}
	QuotationDto struct {
		ID          int64     `json:"id"`
		Code        string    `json:"code"`
		Status      string    `json:"status"`
		Currency    string    `json:"currency"`
		PostingDate time.Time `json:"posting_date"`
		PostingTime string    `json:"posting_time"`
		Tz          string    `json:"tz"`
		ValidTill   time.Time `json:"valid_till"`

		PartyID   int64  `json:"party_id"`
		PartyName string `json:"party_name"`
		PartyUUID string `json:"party_uuid"`
		PartyType string `json:"party_type"`

		PriceListField
		AccountingDimensionDto
	}

	QuotationDetailDto struct {
		Quotation      QuotationDto           `json:"quotation"`
		AcctDimensions AccountingDimensionDto `json:"acc_dimensions"`
	}
)

func QuotationDtoFromModel(m *model.Quotation) QuotationDto {
	return QuotationDto{
		ID:     m.ID,
		Code:   m.Code,
		Status: m.Status,
	}
}
