package dto

import (
	"time"
)

type CreatePurchaseOrderRequest struct {
	AuthParams
	Body struct {
		PartyUUID string    `json:"party_uuid" required:"true"`
		Date      time.Time `json:"date" required:"true"`
		// Name         string         `json:"name" required:"true" minLength:"1" maxLength:"50"`
		DeliveryDate *time.Time     `json:"delivery_date" required:"false"`
		Currency     CurrencyDto    `json:"currency"`
		Lines        []LineItemData `json:"lines" required:"true"`
		PartyType    string         `json:"party_type" required:"true"`
	}
}

// func OrderLineDtoFromModel(m *model.OrderLine) ItemLineDto {
// 	r := ItemLineDto{}
// 	r.Amount = m.Amount
// 	r.Quantity = m.Quantity

// 	if m.ItemPrice.ID != 0 {
// 		r.ItemPrice.UUID = m.ItemPrice.UUID
// 		if m.ItemPrice.Item.ID != 0 {
// 			r.ItemPrice.ItemName = m.ItemPrice.Item.Name
// 			r.ItemPrice.ItemUUID = m.ItemPrice.Item.UUID
// 			r.ItemPrice.ItemCode = m.ItemPrice.Item.Code
// 		}
// 		if m.ItemPrice.UnitOfMeasure.ID != 0 {
// 			r.ItemPrice.Uom = m.ItemPrice.UnitOfMeasure.Code
// 		}
// 	}
// 	return r
// }
