package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CreatePriceListRequest struct {
		AuthParams
		Body PriceListData
		
	}
	EditPriceListRequest struct {
		Body struct {
			ID int64 `json:"id" required:"true"`
			PriceListData
		}
	}

	RequestPriceLists struct {
		AuthParams
		OptionalQueryParams
		PaginationParams
		IsSelling string `query:"is_selling" required:"false"`
		IsBuying  string `query:"is_buying" required:"false"`
	}

	PriceListData struct {
		Name      string `json:"name" required:"true" minLength:"3" maxLength:"50"`
		Currency  string `json:"currency" required:"true" minLength:"2" maxLength:"3"`
		IsBuying  bool   `json:"isBuying" required:"true"`
		IsSelling bool   `json:"isSelling" required:"true"`
	}

	PriceListField struct {
		PriceList *string `json:"price_list"`
		PriceListID *int64 `json:"price_list_id"`
		PriceListUUID *string `json:"price_list_uuid"`
	}

	PriceListDto struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		UUID      string    `json:"uuid"`
		CreatedAt time.Time `json:"created_at"`
		IsBuying  bool      `json:"is_buying"`
		IsSelling bool      `json:"is_selling"`
		Currency  string    `json:"currency"`
	}
)

func PriceListDtoFromModel(m *model.PriceList) PriceListDto {
	r := PriceListDto{}
	r.ID = m.ID
	r.Name = m.Name
	r.UUID = m.UUID
	r.CreatedAt = m.CreatedAt
	r.IsBuying = m.IsBuying
	r.IsSelling = m.IsSelling
	r.Currency = m.Currency
	return r
}
