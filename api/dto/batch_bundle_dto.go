package dto

import "time"

type (
	CreateBatchBundleRequest struct {
		Body struct {
			BatchBundleData
		}
	}

	BatchBundleData struct {
		WarehouseID int64     `json:"warehouse_id"`
		ItemID      int64     `json:"item_id"`
		VoucherType string    `json:"voucher_type"`
		VoucherCode string    `json:"voucher_code"`
		PostingDate time.Time `json:"posting_date"`
		PostingTime string    `json:"Posting_time"`
		SerialNos   []string  `json:"serial_nos"`
	}

	BatchBundleDto struct {
		ID            int64  `json:"id"`
		BatchBundleNo string `json:"batch_bundle_no"`
		VoucherType   string `json:"voucher_type"`
		CreatedAt     string `json:"created_at"`

		Item          string `json:"item"`
		ItemID        int64  `json:"item_id"`
		ItemCode      string `json:"item_code"`
		Warehouse     string `json:"warehouse"`
		WarehouseUUID string `json:"warehouse_uuid"`
	}
)
