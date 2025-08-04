package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CreateStockEntryRequest struct {
		Body StockEntryBody
	}
	EditStockEntryRequest struct {
		Body StockEntryBody
	}

	StockEntryBody struct {
		StockEntry      StockEntryData  `json:"stock_entry" required:"true"`
		CreateItemLines CreateItemLines `json:"items" required:"true"`
	}
	StockEntryData struct {
		ID     int64            `json:"id" required:"false"`
		Fields StockEntryFields `json:"fields" required:"true"`
	}

	StockEntryFields struct {
		EntryType    string    `json:"entry_type"`
		Currency     string    `json:"currency" required:"true"`
		PostingDate  time.Time `json:"posting_date" required:"true"`
		PostingTime  string    `json:"posting_time" required:"true"`
		Tz           string    `json:"tz" required:"true"`
		ProjectID    *int64    `json:"project_id" required:"false"`
		CostCenterID *int64    `json:"cost_center_id" required:"false"`

		SourceWarehouseID *int64 `json:"source_warehouse_id" required:"false"`
		TargetWarehouseID *int64 `json:"target_warehouse_id" required:"false"`
	}

	StockEntryDto struct {
		ID          int64  `json:"id"`
		Code        string `json:"code"`
		Status      string `json:"status"`
		Currency    string `json:"currency"`
		EntryType   string `json:"entry_type"`
		PostingDate string `json:"posting_date"`
		PostingTime string `json:"posting_time"`
		Tz          string `json:"tz"`

		SourceWarehouse     *string `json:"source_warehouse"`
		SourceWarehouseID   *int64  `json:"source_warehouse_id"`
		SourceWarehouseUUID *string `json:"source_warehouse_uuid"`

		TargetWarehouse     *string `json:"target_warehouse"`
		TargetWarehouseID   *int64  `json:"target_warehouse_id"`
		TargetWarehouseUUID *string `json:"target_warehouse_uuid"`
		AccountingDimensionDto
	}
)

func StockEntryDtoFromModel(m *model.StockEntry) StockEntryDto {
	return StockEntryDto{
		ID:     m.ID,
		Code:   m.Code,
		Status: m.Status,
	}
}
