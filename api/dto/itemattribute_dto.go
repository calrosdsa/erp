package dto

import (
	"erp/gen/db/model"
	"time"
)

type CreateItemAttributeRequest struct {
	AuthParams
	Body struct {
		Name   string                  `json:"name"`
		Values []ItemAttributeValueDto `json:"values"`
	}
}

type ItemAttributeValueDto struct {
	ID              int32  `json:"id" required:"false"`
	ItemAttributeID int64  `json:"itemAttributeId" required:"false"`
	Ordinal         int32  `json:"ordinal"`
	Value           string `json:"value"`
	Abbreviation    string `json:"abbreviation"`
}

type UpsertItemAttributeValueRequest struct {
	AuthParams
	Body struct {
		ItemAttributeValueDto
	}
}

type ItemAttributeDto struct {
	ID int64 `json:"id"`
	UUID                string                  `json:"uuid"`
	Name                string                  `json:"name"`
	CreatedAt           time.Time               `json:"created_at"`
	ItemAttributeValues []ItemAttributeValueDto `json:"item_attribute_values"`
}

func ItemAttributeDtoFromModel(m *model.ItemAttribute) ItemAttributeDto {
	r := ItemAttributeDto{}
	r.ID = m.ID
	r.Name = m.Name
	r.UUID = m.UUID
	r.CreatedAt = m.CreatedAt
	return r
}


func ItemAttributeValueDtoFromModel(m *model.ItemAttributeValue) ItemAttributeValueDto {
	r := ItemAttributeValueDto{}
	r.Abbreviation = m.Abbreviation
	r.Value = m.Value
	r.ID = m.ID
	r.ItemAttributeID = m.ItemAttributeID
	r.Ordinal = m.Ordinal
	return r
}