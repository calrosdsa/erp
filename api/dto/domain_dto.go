package dto

import "erp/gen/db/model"

type CurrencyDto struct {
	Code string `json:"code"`
}

func CurrencyDtoFromModel(m *model.Currency) CurrencyDto {
	r :=  CurrencyDto{}
	r.Code = m.Code
	return r
}
