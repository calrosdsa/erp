package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	ItemPricesRequest struct {
		DefaultListParams
		Currency    string `query:"currency"`
		PriceListID string `query:"price_list_id"`
		ItemID      string `query:"item_id"`
	}

	ItemPriceRequestData struct {
		AuthParams
		Body ItemPriceData
	}

	EditItemPriceRequest struct {
		Body struct {
			ID int64 `json:"id" required:"true"`
			ItemPriceData
		}
	}

	ItemPriceData struct {
		ID     int64           `json:"id" required:"false"`
		Fields ItemPriceFields `json:"fields"`
		Action string          `json:"action" required:"false"`
	}

	ItemPriceFields struct {
		ItemID          int64 `json:"item_id"`
		PriceListID     int64 `json:"price_list_id"`
		ItemQuantity    int32 `json:"item_quantity"`
		Rate            int64 `json:"rate"`
		UnitOfMeasureID int64 `json:"unit_of_measure_id"`
	}

	RequestItemPricesForOrder struct {
		OptionalQueryParams
		AuthParams
		Currency    string `query:"currency" required:"true"`
		IsBuying    bool   `query:"isBuying" required:"false"`
		IsSelling   bool   `query:"isSelling" required:"false"`
		Enabled     string `query:"enabled" required:"false"`
		PriceListID string `query:"price_list_id" required:"false"`
	}

	ItemPriceDetailDto struct {
		UUID         string     `json:"uuid"`
		Code         string     `json:"code"`
		ItemQuantity int32      `json:"item_quantity" required:"true"`
		Rate         int32      `json:"rate" required:"true"`
		ValidFrom    *time.Time `json:"valid_from" required:"false"`
		ValidUpTo    *time.Time `json:"valid_up_to" required:"false"`
		CreatedAt    time.Time  `json:"created_at"`
		//No part of  the model
		//Item
		ItemName string `json:"item_name"`
		ItemUUID string `json:"item_uuid"`
		ItemID   int64  `json:"item_id"`
		ItemCode string `json:"item_code"`
		//Uom
		Uom string `json:"uom"`
		//PriceList
		Currency string `json:"currency"`
	}
	ItemPriceDto struct {
		ID           int64     `json:"id"`
		UUID         string    `json:"uuid"`
		Code         string    `json:"code"`
		ItemQuantity int32     `json:"item_quantity" required:"true"`
		Rate         int64     `json:"rate" required:"true"`
		CreatedAt    time.Time `json:"created_at"`
		//No part of  the model
		//Item
		ItemName string `json:"item_name"`
		ItemUUID string `json:"item_uuid"`
		ItemCode string `json:"item_code"`
		ItemID   int64  `json:"item_id"`
		//Item Price
		ItemPriceUomID int64  `json:"item_price_uom_id"`
		ItemPriceUom   string `json:"item_price_uom"`
		//Uom
		Uom   string `json:"item_uom"`
		UomID int64  `json:"item_uom_id"`
		//PriceList
		PriceListID       int64  `json:"price_list_id"`
		PriceListName     string `json:"price_list_name"`
		PriceListUUID     string `json:"price_list_uuid"`
		PriceListCurrency string `json:"price_list_currency"`
	}
)

func ItemPriceDtoFromModel(m *model.ItemPrice) ItemPriceDto {
	r := ItemPriceDto{}
	r.ID = m.ID
	r.UUID = m.UUID
	r.Rate = m.Rate
	r.ItemQuantity = m.ItemQuantity
	r.CreatedAt = m.CreatedAt

	return r
}
