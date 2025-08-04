package dto

import "erp/gen/db/model"

type (
	BankAccountsRequest struct {
		DefaultListParams
		AccountName      string `query:"account_name" required:"false"`
		PartyID          string `query:"party_id" required:"false"`
		IsCompanyAccount string   `query:"is_company_account" required:"false"`

		CreatedAt string `query:"created_at" required:"false"`
		UpdatedAt string `query:"updated_at" required:"false"`
	}
	BankAccountDataRequest struct {
		Body BankAccountData
	}
	BankAccountData struct {
		ID     int64             `json:"id"`
		Fields BankAccountFields `json:"fields"`
	}

	BankAccountFields struct {
		AccountName     string `json:"account_name" required:"true"`
		BankAccountType string `json:"bank_account_type" required:"true"`

		BankID  int64  `json:"bank_id"`
		PartyID *int64 `json:"party_id"`

		Iban              *string `json:"iban" required:"false"`
		BranchCode        *string `json:"branch_code" required:"false"`
		BankAccountNumber *string `json:"bank_account_number" required:"false"`

		IsCompanyAccount bool   `json:"is_company_account"`
		CompanyAccountID *int64 `json:"company_account_id" required:"false"`
	}

	BankAccountDto struct {
		ID     int64  `json:"id"`
		UUID   string `json:"uuid"`
		Status string `json:"status"`

		AccountName     string `json:"account_name"`
		BankAccountType string `json:"bank_account_type"`

		BankAccountNumber *string `json:"bank_account_number"`
		Iban              *string `json:"iban"`
		BranchCode        *string `json:"branch_code"`

		Party     *string `json:"party"`
		PartyID   *int64  `json:"party_id"`
		PartyUUID *string `json:"party_uuid"`
		PartyType *string `json:"party_type"`

		Bank     string `json:"bank"`
		BankID   int64  `json:"bank_id"`
		BankUUID string `json:"bank_uuid"`

		IsCompanyAccount bool `json:"is_comapny_account"`

		CompanyAccount     *string `json:"company_account"`
		CompanyAccountID   *int64  `json:"company_account_id"`
		CompanyAccountUUID *string `json:"company_account_uuid"`
		CompanyAccountCurrency *string `json:"company_account_currency"`
	}
)

func BankAccountFromModel(m model.BankAccount) BankAccountDto {
	r := BankAccountDto{}
	r.ID = m.ID
	r.UUID = m.UUID
	r.AccountName = m.AccountName
	return r
}
