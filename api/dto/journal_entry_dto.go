package dto

import (
	"erp/gen/db/model"
	"time"
)

type (

	

	JournalEntryRequestData struct {
		Body JournalEntryData
	}
	JournalEntryData struct {
		ID int64 `json:"id" required:"false"`
		Fields JournalEntryFields `json:"fields"`
		EntryLines  []JournalEntryLineData `json:"entry_lines"`
	}
	JournalEntryFields struct {
		PostingDate time.Time `json:"posting_date"`
		EntryType string `json:"entry_type"`
	}


	JournalEntryLineData struct {
		Debit        float64 `json:"debit" required:"true"`
		Credit       float64 `json:"credit" required:"true"`
		LedgerID     int64   `json:"ledger_id" required:"true"`
		CostCenterID *int64  `json:"cost_center_id"`
		ProjectID    *int64  `json:"project_id"`
	}

	JournalEntryDetailDto struct {
		JournalEntry      JournalEntryDto       `json:"journal_entry"`
		JournalEntryLines []JournalEntryLineDto `json:"journal_entry_lines"`
	}

	JournalEntryLineDto struct {
		ID       int32  `json:"id"`
		Debit    int    `json:"debit"`
		Credit   int    `json:"credit"`
		Currency string `json:"currency"`

		Account   string `json:"account"`
		AccountID int64  `json:"account_id"`

		CostCenter   *string `json:"cost_center"`
		CostCenterID *int64  `json:"cost_center_id"`

		Project   *string `json:"project"`
		ProjectID *int64  `json:"project_id"`
	}

	JournalEntryDto struct {
		ID          int64  `json:"id"`
		Code        string `json:"code"`
		Status      string `json:"status"`
		EntryType   string `json:"entry_type"`
		PostingDate string `json:"posting_date"`
		CreatedAt   string `json:"created_at"`
	}
)

func JournalEntryDtoFromModel(m *model.JournalEntry) JournalEntryDto {
	return JournalEntryDto{
		Code:   m.Code,
		ID:     m.ID,
		Status: m.Status,
	}
}
