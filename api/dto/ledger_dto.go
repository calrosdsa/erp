package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	LedgerDataRequest struct {
		Body LedgerData
	}

	LedgerData struct {
		ID     int64        `json:"id" required:"false"`
		Fields LedgerFields `json:"fields"`
	}
	LedgerFields struct {
		LedgerParent    *int64  `json:"ledger_parent" required:"false"`
		AccountType     *string `json:"account_type" required:"false"`
		Name            string  `json:"name"`
		IsGroup         bool    `json:"is_group"`
		LedgerNo        *string `json:"ledger_no" required:"false"`
		AccountRootType string  `json:"account_root_type"`
		ReportType      *string `json:"report_type" required:"false"`
		CashFlowSection *string `json:"cash_flow_section" required:"false"`
		IsOffsetAccount bool    `json:"is_offset_account" required:"false"`
	}

	LedgerDto struct {
		ID              int64     `json:"id"`
		UUID            string    `json:"uuid"`
		Name            string    `json:"name"`
		Description     string    `json:"description"`
		IsGroup         bool      `json:"is_group"`
		Status          string    `json:"status"`
		AccountType     string    `json:"account_type"`
		AccountRootType string    `json:"account_root_type"`
		ReportType      string    `json:"report_type"`
		CashFlowSection string    `json:"cash_flow_section"`
		LedgerNo        *string   `json:"ledger_no"`
		CreatedAt       time.Time `json:"created_at"`
		IsOffsetAccount bool      `json:"is_offset_account"`

		LedgerAccountDto
	}

	LedgerAccountDto struct {
		Currency  string `json:"currency"`
		CanCredit bool   `json:"can_credit"`
		CanDebit  bool   `json:"can_debit"`
		Limit     int    `json:"limit"`
	}

	LedgerDetailDto struct {
		LedgerDto
		// LedgerAccountDto

		//Parent
		Parent     *string `json:"parent"`
		ParentID   *int64  `json:"parent_id"`
		ParentUUID *string `json:"parent_uuid"`
	}

	LedgersRequest struct {
		DefaultListParams
		Name            string `query:"name" required:"false"`
		IsCreditBalance string `query:"is_credit_balance" required:"false"`
		IsDebitBalance  string `query:"is_debit_balance" required:"false"`
		IsGroup         string `query:"is_group" required:"false"`
		AccountType     string `query:"account_type" required:"false"`
	}
)

func LedgerDtoFromModel(ledger *model.Ledger) LedgerDto {
	r := LedgerDto{}
	r.ID = ledger.ID
	r.Name = ledger.Name
	r.UUID = ledger.UUID
	return r
}
