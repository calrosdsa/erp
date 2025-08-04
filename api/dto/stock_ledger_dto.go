package dto

import "time"

type (
	RequestStockLedger struct {
		AuthParams
		FromDate    string `query:"from_date" required:"true"`
		ToDate      string `query:"to_date" required:"true"`
		VoucherNo   string `query:"voucher_no"`
		ItemID      string `query:"item_id"`
		WarehouseID string `query:"warehouse_id"`
	}

	RequestStockBalance struct {
		AuthParams
		FromDate    string `query:"from_date"`
		ToDate      string `query:"to_date"`
		ItemID      string `query:"item_id"`
		WarehouseID string `query:"warehouse_id"`
	}

	StockLedgerEntryDto struct {
		Date          time.Time `json:"date"`
		ItemName      string    `json:"item_name"`
		ItemUUID      string    `json:"item_uuid"`
		ItemID        int64     `json:"item_id"`
		ItemCode      string    `json:"item_code"`
		ItemGroupName string    `json:"item_group_name"`
		ItemGroupUUID string    `json:"item_group_uuid"`
		StockUom      string    `json:"stock_uom"`
		InQty         int32     `json:"in_qty"`
		OutQty        int32     `json:"out_qty"`
		BalanceQty    int32     `json:"balance_qty"`
		WarehouseName string    `json:"warehouse_name"`
		WarehouseUUID string    `json:"warehouse_uuid"`
		IncomingRate  int32     `json:"incoming_rate"`
		AverageRate   int32     `json:"average_rate"`
		ValuationRate int32     `json:"valuation_rate"`
		BalanceValue  int32     `json:"balance_value"`
		VoucherType   string    `json:"voucher_type"`
		VoucherNo     string    `json:"voucher_no"`
		Currency      string    `json:"currency"`
	}

	StockBalanceEntryDto struct {
		Date          time.Time `json:"date"`
		ItemName      string    `json:"item_name"`
		ItemID        int64     `json:"item_id"`
		ItemUUID      string    `json:"item_uuid"`
		ItemCode      string    `json:"item_code"`
		ItemGroupName string    `json:"item_group_name"`
		ItemGroupUUID string    `json:"item_group_uuid"`
		StockUom      string    `json:"stock_uom"`
		InQty         int32     `json:"in_qty"`
		OutQty        int32     `json:"out_qty"`
		BalanceQty    int32     `json:"balance_qty"`
		BalanceValue  int32     `json:"balance_value"`
		WarehouseName string    `json:"warehouse_name"`
		WarehouseUUID string    `json:"warehouse_uuid"`
		AverageRate   int32     `json:"average_rate"`
		Currency      string    `json:"currency"`
	}
)
