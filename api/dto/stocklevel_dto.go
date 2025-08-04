package dto

import (
	"erp/gen/db/model"
	"time"
)

type StockLevelDto struct {
	UUID                string    `json:"uuid"`
	CreatedAt           time.Time `json:"created_at"`
	Enabled             bool      `json:"enabled"`
	Stock               int32     `json:"stock"`
	OutOfStockThreshold int32     `json:"out_of_stock_threshold"`
	// No segment of the entity
	WarehouseName string `json:"warehouse_name"`
	WarehouseUUID string `json:"warehouse_uuid"`

	ItemName string `json:"item_name"`
	ItemUUID string `json:"item_uuid"`
}

func StockLevelDtoFromModel(m *model.StockLevel) StockLevelDto {
	r := StockLevelDto{}
	r.UUID = m.UUID
	r.CreatedAt = m.CreatedAt
	r.Enabled = m.Enabled
	r.Stock = m.Stock
	r.OutOfStockThreshold = m.OutOfStockThreshold
	return r
}
