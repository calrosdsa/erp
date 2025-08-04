package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CreateCompanyRequest struct {
		AuthParams
		Body struct {
			CompanyEditableFields
		}
	}

	CompanyEditableFields struct {
		Name string `json:"name" required:"true" minLength:"3" maxLength:"50"`
		// Name string `json:"name" required:"true" minLength:'3" maxLength:"50"`
		ParentID *int64 `json:"parentId" required:"false"`
		// Logo string `json:"logo" required:"false"`
		// SiteUrl string `json:"site"`
	}

	AccountSettingDataRequest struct {
		Body AccountSettingData
	}

	AccountSettingData struct {
		ID     int64                `json:"id"`
		Fields AccountSettingFields `json:"fields"`
	}

	AccountSettingFields struct {
		BankAccount           int64 `json:"bank_account" required:"true"`
		CashAccunt            int64 `json:"cash_accunt" required:"true"`
		PayableAccount        int64 `json:"payable_account" required:"true"`
		CostOfGoodSoldAccount int64 `json:"cost_of_good_sold_account" required:"true"`
		ReceivableAccount     int64 `json:"receivable_account" required:"true"`
		IncomeAccount         int64 `json:"income_account" required:"true"`
	}

	CompanyParentDto struct {
		ID       int64  `json:"id"`
		ParentID *int64 `json:"parent_id"`
	}

	CompanyDto struct {
		ID        int64     `json:"id"`
		UUID      string    `json:"uuid"`
		CreatedAt time.Time `json:"created_at"`
		Name      string    `json:"name"`
		Ordinal   int       `json:"ordinal"`
		Logo      *string   `json:"logo"`
		SiteURL   *string   `json:"site_url"`
	}

	AccountSettingsDto struct {
		CashAcct     string `json:"cash_acct"`
		CashAcctID   int64  `json:"cash_acct_id"`
		CashAcctUUID string `json:"cash_acct_uuid"`

		BankAcct     string `json:"bank_acct"`
		BankAcctID   int64  `json:"bank_acct_id"`
		BankAcctUUID string `json:"bank_acct_uuid"`

		PayableAcct     string `json:"payable_acct"`
		PayableAcctID   int64  `json:"payable_acct_id"`
		PayableAcctUUID string `json:"payable_acct_uuid"`

		CostOfGoodsSoldAcct     string `json:"cost_of_goods_sold_acct"`
		CostOfGoodsSoldAcctID   int64  `json:"cost_of_goods_sold_acct_id"`
		CostOfGoodsSoldAcctUUID string `json:"cost_of_goods_sold_acct_uuid"`

		ReceivableAcct     string `json:"receivable_acct"`
		ReceivableAcctID   int64 `json:"receivable_acct_id"`
		ReceivableAcctUUID string `json:"receivable_acct_uuid"`

		IncomeAcct     string `json:"income_acct"`
		IncomeAcctID   int64  `json:"income_acct_id"`
		IncomeAcctUUID string `json:"income_acct_uuid"`
	}

	CompanyDefaultsDto struct {
		Currency string `json:"currency"`
	}
)

func CompanyDefaultsDtoFromModel(m *model.CompanyDefault) CompanyDefaultsDto {
	return CompanyDefaultsDto{
		Currency: m.Currency,
	}
}

func CompanyDTOFromModel(m *model.Company) CompanyDto {
	p := CompanyDto{}
	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.Name = m.Name
	p.UUID = m.UUID
	p.Ordinal = int(m.Ordinal)
	p.Logo = m.Logo
	p.SiteURL = m.SiteURL
	return p
}

func CompanyDTOFromModelWithID(m *model.Company) CompanyDto {
	p := CompanyDto{}
	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.Name = m.Name
	p.UUID = m.UUID
	p.Ordinal = int(m.Ordinal)
	p.Logo = m.Logo
	p.SiteURL = m.SiteURL
	return p
}

//  CompanyDetailDto struct {
// 	UUID      string      `json:"string"`
// 	CreatedAt time.Time   `json:"created_at"`
// 	Name      string      `json:"name"`
// 	Ordinal   string      `json:"ordinal"`
// 	Code      *string     `json:"code"`
// 	Logo      *string     `json:"logo"`
// 	SiteURL   *string     `json:"site_url"`
// 	Parent    *CompanyDto `json:"parent"`
// }
