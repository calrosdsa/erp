package dto

import "erp/gen/db/model"

type (
	CreateChargesTemplateRequest struct {
		Body ChargesTemplateBody
	}
	ChargesTemplateBody struct {
		ChargesTemplate     ChargesTemplateData `json:"charges_template"`
		CreateTaxAndCharges CreateTaxAndChanges `json:"tax_and_charges"`
	}

	EditChargesTemplateRequest struct {
		Body struct {
			ID int64 `json:"id" required:"true"`
			ChargesTemplateData
		}
	}

	ChargesTemplateData struct {
		Name string `json:"name" required:"true"`
	}

	ChargesTemplateDto struct {
		ID        int64  `json:"id"`
		UUID      string `json:"uuid"`
		Name      string `json:"name"`
		Status    string `json:"status"`
		CreatedAt string `json:"created_at"`
	}
)

func ChargesTemplateFromModel(m *model.ChargesTemplate) ChargesTemplateDto {
	return ChargesTemplateDto{
		ID: m.ID,
		Name: m.Name,
		UUID: m.UUID,
	}
}
