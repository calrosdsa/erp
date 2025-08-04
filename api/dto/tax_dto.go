package dto

import (
	"erp/gen/db/model"
	"time"
)

type CreateTaxRequest struct {
	AuthParams
	Body struct {
		TaxEditableFields
	}
}

type TaxEditableFields struct {
	Name    string  `json:"name" required:"true" minLength:"3" maxLength:"50"`
	Value   float64 `json:"value" required:"true"`
	Enabled bool    `json:"enabled" required:"true"`
}

type TaxDto struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	UUID      string    `json:"uuid"`
	Enabled   bool      `json:"enabled"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}



func TaxDtoFromModel(tax *model.Tax) TaxDto {
	r := TaxDto{}
	r.ID = tax.ID
	r.Name = tax.Name
	r.UUID = tax.UUID
	r.Value = tax.Value
	r.Enabled = tax.Enabled
	r.CreatedAt = tax.CreatedAt
	return r
}
