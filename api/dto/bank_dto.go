package dto

import "erp/gen/db/model"

type (
	BanksRequest struct {
		DefaultListParams
		Name      string `query:"name" required:"false"`
		CreatedAt string `query:"created_at" required:"false"`
		UpdatedAt string `query:"updated_at" required:"false"`
	}
	BankDataRequest struct {
		Body BankData
	}
	BankData struct {
		ID     int64      `json:"id"`
		Fields BankFields `json:"fields"`
	}

	BankFields struct {
		Name string `json:"name"`
	}

	BankDto struct {
		ID     int64  `json:"id"`
		UUID   string `json:"uuid"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
)

func BankFromModel(m model.Bank) BankDto {
	r := BankDto{}
	r.ID = m.ID
	r.UUID = m.UUID
	r.Name = m.Name
	return r
}
