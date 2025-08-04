package dto

import "erp/gen/db/model"

type (
	CreateCurrencyExchangeRequest struct {
		Body struct {
			CurrencyExchangeData
		}
	}
	EditCurrencyExchangeRequest struct {
		Body struct {
			ID int64 `json:"id"`
			CurrencyExchangeData
		}
	}
	CurrencyExchangeData struct {
		Name         string `json:"name" required:"true"`
		FromCurrency string `json:"from_currency" required:"true"`
		ToCurrency   string `json:"to_currency" required:"true"`
		ExchangeRate float64  `json:"exchange_rate" required:"true"`
		ForBuying    bool   `json:"for_buying" required:"true"`
		ForSelling   bool   `json:"for_selling" required:"true"`
	}

	CurrencyExchangeDto struct {
		ID           int64  `json:"id"`
		UUID         string `json:"uuid"`
		Name         string `json:"name"`
		Status       string `json:"status"`
		FromCurrency string `json:"from_currency"`
		ToCurrency   string `json:"to_currency"`
		ExchangeRate int32  `json:"exchange_rate"`
		ForBuying    bool   `json:"for_buying"`
		ForSelling   bool   `json:"for_selling"`
	}
)

func CurrencyExchangeDtoFromModel(m *model.CurrencyExchange) CurrencyExchangeDto {
	return CurrencyExchangeDto{
		ID:m.ID,
		UUID:m.UUID,
		Name:m.Name,
	}
}
