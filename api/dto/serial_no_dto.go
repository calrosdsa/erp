package dto

type (
	RequestSerialNos struct {
		PaginationParams
		OptionalQueryParams
		SerialNo string `query:"serial_no" required:"false"`
	}
	RequestSerialNoTransactions struct {
		FromDate      string `query:"from_date" required:"true"`
		ToDate        string `query:"to_date" required:"true"`
		VoucherCode   string `query:"voucher_code" required:"false"`
		SerialNo      string `query:"serial_no" required:"false"`
		BatchBundleNo string `query:"batch_bundle_no" required:"false"`
		ItemID        string `query:"item_id" required:"false"`
		WarehouseID   string `query:"warehouse_id" required:"false"`
	}
	SerialNoDto struct {
		ID       int64  `json:"id"`
		SerialNo string `json:"serial_no"`

		ItemName  string `json:"item_name"`
		ItemID    int64  `json:"item_id"`
		ItemCode  string `json:"item_code"`
		CreatedAt string `json:"created_at"`
		Status    string `json:"status"`

		Warehouse     string `json:"warehouse"`
		WarehouseUUID string `json:"warehouse_uuid"`
	}

	SerialNoTransactionDto struct {
		ID            int64  `json:"id"`
		SerialNo      string `json:"serial_no"`
		ValuationRate int32  `json:"valuation_rate"`
		Qty           int32  `json:"qty"`
		Status        string `json:"status"`

		Item     string `json:"item_name"`
		ItemID   int64  `json:"item_id"`
		ItemCode string `json:"item_code"`

		Warehouse     string `json:"warehouse_name"`
		WarehouseID   int64  `json:"warehouse_id"`
		WarehouseUUID string `json:"warehouse_uuid"`

		PostingDate string `json:"posting_date"`
		PostingTime string `json:"posting_time"`
		VoucherType string `json:"voucher_type"`
		VoucherCode string `json:"voucher_code"`

		BatchBundleNo string `json:"batch_bundle_no"`
	}
)
