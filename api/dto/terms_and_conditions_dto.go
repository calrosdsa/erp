package dto

import "erp/gen/db/model"

type (
	TermsAndConditionsRequest struct {
		DefaultListParams
		Name string `query:"name"`
		CreatedAt string `query:"created_at"`
	}

	TermsAndConditionsDataRequest struct {
		Body TermsAndConditionsData
	}
	TermsAndConditionsData struct {
		ID     int64                   `json:"id" required:"false"`
		Fields TermsAndConditionFields `json:"filds" required:"true"`
	}
	TermsAndConditionFields struct {
		Name             string `json:"name" required:"true"`
		TermsAndConditions string `json:"terms_and_conditions" required:"true"`
	}
	TermsAndConditionsDto struct {
		ID               int64  `json:"id"`
		UUID string `json:"uuid"`
		Name             string `json:"name"`
		TermsAndConditions string `json:"terms_and_conditions"`
		Status           string `json:"status"`
	}
)

func TermsAndConditionFromModel(m model.TermsAndCondition)TermsAndConditionsDto {
	return TermsAndConditionsDto{
		ID:m.ID,
		UUID: m.UUID,
		Name:m.Name,
		TermsAndConditions: m.TermsAndConditions,
		Status: m.Status,
	}
}