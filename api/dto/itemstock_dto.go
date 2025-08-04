package dto

type AddStockLevelRequest struct {
	AuthParams
	Body struct {
		Stock               int32  `json:"stock" required:"true"`
		Enabled             bool   `json:"enabled" required:"true"`
		OutOfStockThreshold int32  `json:"outOfStockThreshold" required:"true"`
		ItemUUID            string `json:"item_uuid" required:"true"`
		WareHouseUUID       string  `json:"warehouse_uuid" required:"true"`
	}
}
