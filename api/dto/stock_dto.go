package dto

import (
	"erp/gen/db/model"
)

type UpsertPriceListRequest struct {
	AuthParams
	Body struct {
		ItemPriceList model.PriceList `json:"itemPriceList"`
	}
}

type UpsertItemPriceRequest struct {
	AuthParams
	Body struct {
		ItemPrice model.ItemPrice `json:"itemPrice"`
	}
}


type PluginDto struct {
	CompanyId uint   `json:"companyId" required:"false"`
	Plugin    string `json:"plugin" required:"true"`
}

type RequestItemPriceByCode struct {
	PaginationParams
	AuthParams
	ItemID string `path:"item_id" required:"true"`
}
