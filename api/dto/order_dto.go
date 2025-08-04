package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	RequestOrders struct {
		PartyType string `path:"party" required:"true"`
		PaginationParams
		OptionalQueryParams
		Code         string `query:"code" required:"false"`
		DeliveryDate string `query:"delivery_date" required:"false"`
		PostingDate  string `query:"posting_date" required:"false"`
		PartyID      string `query:"party_id" required:"false"`
		ID           string `query:"id" required:"false"`
		PricingID    string `query:"pricing_id" required:"false"`
		ProjectID    string `query:"project_id" required:"false"`
		CostCenterID string `query:"cost_center_id" required:"false"`
	}

	CreateOrderRequest struct {
		Body OrderBody
	}

	EditOrderRequest struct {
		Body OrderBody
	}

	OrderBody struct {
		Order               OrderData           `json:"order" required:"true"`
		CreateItemLines     CreateItemLines     `json:"items" required:"true"`
		CreateTaxAndCharges CreateTaxAndChanges `json:"tax_and_charges"`
	}

	OrderFields struct {
		PartyID      int64      `json:"party_id" required:"true"`
		PostingDate  time.Time  `json:"posting_date" required:"true"`
		PostingTime  string     `json:"posting_time" required:"true"`
		Tz           string     `json:"tz"`
		Currency     string     `json:"currency"`
		DeliveryDate *time.Time `json:"delivery_date" required:"false"`
		ProjectID    *int64     `json:"project_id" required:"false"`
		CostCenterID *int64     `json:"cost_center_id" required:"false"`
		PriceListID  *int64     `json:"price_list_id" required:"false"`
	}

	OrderData struct {
		ID             int64       `json:"id" required:"false"`
		OrderPartyType string      `json:"order_party_type" required:"true"`
		Fields         OrderFields `json:"fields"`
		References     []*int64    `json:"references" required:"false"`
		TotalAmount    float64     `json:"total_amount"`
	}

	OrderDto struct {
		ID           int64      `json:"id"`
		CreatedAt    time.Time  `json:"created_at"`
		PostingDate  time.Time  `json:"posting_date"`
		PostingTime  string     `json:"posting_time"`
		Tz           string     `json:"tz"`
		DeliveryDate *time.Time `json:"delivery_date"`
		Code         string     `json:"code"`
		Currency     string     `json:"currency"`
		Status       string     `json:"status"`

		PartyName string `json:"party_name"`
		PartyID   int64  `json:"party_id"`
		PartyUuid string `json:"party_uuid"`
		PartyType string `json:"party_type"`

		PriceListField

		AccountingDimensionDto
		//ProgressOrder
		TotalItems    int32 `json:"total_items"`
		ReceivedItems int32 `json:"received_items"`
		TotalAmount   int32 `json:"total_amount"`
		BilledAmount  int32 `json:"billed_amount"`
		//Optional attibutes
	}

	OrderDetailDto struct {
		Order          OrderDto               `json:"order"`
		AcctDimensions AccountingDimensionDto `json:"acc_dimensions"`
		PartyAddress   AddressDto             `json:"party_address"`
		PartyContact   ContactDto             `json:"party_contact"`
		CompanyAddress AddressDto             `json:"company_addresss"`
	}
)

func OrderDtoFromModel(m *model.Order) OrderDto {
	r := OrderDto{}
	// r.Name = m.Name
	r.CreatedAt = m.CreatedAt
	r.DeliveryDate = m.DeliveryDate
	r.Code = m.Code
	r.PostingDate = m.PostingDate
	r.PostingTime = m.PostingTime
	r.Status = m.Status
	r.Currency = m.Currency
	return r
}
