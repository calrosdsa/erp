package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CreateWareHouseRequest struct {
		AuthParams
		Body struct {
			WarehouseData
		}
	}
	EditWarehouseRequest struct {
		Body struct {
			ID int64 `json:"id" required:"true"`
			WarehouseData
		}
	}

	WarehouseOptionalField struct {
		Warehouse     *string `json:"warehouse"`
		WarehouseID   *int64  `json:"warehouse_id"`
		WarehouseUUID *string `json:"warehouse_uuid"`
	}

	WarehouseField struct {
		Warehouse     string `json:"warehouse"`
		WarehouseID   int64  `json:"warehouse_id"`
		WarehouseUUID string `json:"warehouse_uuid"`
	}

	WarehouseData struct {
		Name     string `json:"name" required:"true"`
		ParentID int64  `json:"parent_id" required:"false"`
		IsGroup  bool   `json:"is_group" required:"true"`
	}

	WareHouseDto struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		UUID      string    `json:"uuid"`
		Enabled   bool      `json:"enabled"`
		IsGroup   bool      `json:"is_group"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func WarehouseDtoFromModel(m *model.WareHouse) WareHouseDto {
	r := WareHouseDto{}
	r.ID = m.ID
	r.UUID = m.UUID
	r.CreatedAt = m.CreatedAt
	r.Name = m.Name
	r.Enabled = m.Enabled
	return r
}
