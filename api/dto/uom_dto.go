package dto

import (
	"erp/gen/db/model"
)

type UOMsRequest struct {
	AuthParams
	OptionalQueryParams
}

type UOMsResponse struct {
	Body struct {
		Results []UOMDto `json:"results"`
	}
}

type UOMDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func UOMFromModel(m *model.UnitOfMeasure) UOMDto {
	r := UOMDto{}
	if m.UnitOfMeasureTranslation.ID != 0 {
		r.Name = m.UnitOfMeasureTranslation.Name
	}
	r.Code = m.Code
	r.ID = m.ID
	return r
}
