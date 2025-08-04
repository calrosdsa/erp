package dto

import "time"

type CreateItemVariantRequest struct {
	AuthParams
	Body struct {
		Name                  string `json:"name" required:"true"`
		ItemUUID              string `json:"item_uuid" required:"true"`
		AttributeValueValueID int32  `json:"attribute_value_id" required:"true"`
	}
}

type ItemVariantDto struct {
	Name      string    `json:"name"`
	UUID      string    `json:"uuid"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`

	AttributeValue        string `json:"attibute_value"`
	AttributeAbbreviation string `json:"attibute_abbreviation"`
}
