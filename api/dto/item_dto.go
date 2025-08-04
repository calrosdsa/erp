package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	ItemRequestData struct {
		AuthParams
		Body ItemData
	}

	EditItemRequest struct {
		Body struct {
			ID int64 `json:"id"`
			ItemData
		}
	}

	CreateItemData struct {
		Item           ItemData            `json:"item"`
		ItemPriceLines []ItemPriceData     `json:"item_price_lines"`
		ItemInventory  ItemInventoryFields `json:"item_inventory"`
	}

	ItemPriceLine struct {
		ID           int64   `json:"id" required:"false"`
		Rate         float64 `json:"rate" required:"true"`
		ItemQuantity int32   `json:"item_quantity" required:"true"`
		PriceListID  int64   `json:"price_list_id" required:"true"`
		UomID        int64   `json:"uom_id" required:"true"`
		Action       string  `json:"action" required:"false"`
	}

	// StockEntryLine struct

	ItemData struct {
		ID            int64               `json:"id" required:"false"`
		Fields        ItemFields          `json:"fields"`
		ItemPrices    []ItemPriceData     `json:"item_price_lines"`
		ItemInventory ItemInventoryFields `json:"item_inventory"`
	}
	ItemFields struct {
		Name            string  `json:"name"`
		Code            *string `json:"code" required:"false"`
		GroupID         *int64  `json:"group_id" required:"false"`
		UnitOfMeasureID int64   `json:"unit_of_measure_id" required:"true"`
		MaintainStock   bool    `json:"maintain_stock" required:"true"`
		Description     *string `json:"description" required:"false"`
	}

	UpdateItemRequest struct {
		AuthParams
		OptionalQueryParams
		Body struct {
			Name     string `json:"name" required:"true" minLength:"1" maxLength:"50"`
			ItemType string `json:"item_type" required:"true" minLength:"1" maxLength:"50"`
		}
	}

	ItemDetailDto struct {
		ID            int64     `json:"id"`
		Name          string    `json:"name"`
		UUID          string    `json:"uuid"`
		Code          *string   `json:"code"`
		CreatedAt     time.Time `json:"created_at"`
		ItemType      string    `json:"item_type"`
		Description   *string   `json:"description"`
		MaintainStock bool      `json:"maintain_stock"`
		Status        string    `json:"status"`

		// Group GroupDto `json:"group"`
		GroupID   *int64  `json:"group_id"`
		GroupName *string `json:"group_name"`
		GroupUuid *string `json:"group_uuid"`

		UomName string `json:"uom_name"`
		UomCode string `json:"uom_code"`
		UomID   int64  `json:"uom_id"`

		ItemInventoryDto

		ItemPrices []ItemPriceDto `json:"item_prices" gorm:"-"`
	}

	ItemDto struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		UUID      string    `json:"uuid"`
		Code      *string   `json:"code"`
		CreatedAt time.Time `json:"created_at"`
		ItemType  string    `json:"item_type"`
		Status    string    `json:"status"`
		UomID     int64     `json:"uom_id"`

		// Uom string `json:"uom"`
	}
)

func ItemDtoFromModel(m *model.Item) ItemDto {
	r := ItemDto{}
	r.ID = m.ID
	r.Name = m.Name
	r.UUID = m.UUID
	r.ItemType = m.ItemType
	r.CreatedAt = m.CreatedAt
	return r
}

func ItemDetailDtoFromModel(m *model.Item) ItemDetailDto {
	r := ItemDetailDto{}
	r.ID = m.ID
	r.Name = m.Name
	r.UUID = m.UUID
	r.Code = m.Code
	r.ItemType = m.ItemType
	r.CreatedAt = m.CreatedAt

	return r
}
