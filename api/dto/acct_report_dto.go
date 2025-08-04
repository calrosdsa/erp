package dto

import "time"

type (
	RequestAccountBalance struct {
		ID string `query:"id" required:"id"`
	}

	RequestGeneralLedger struct {
		AuthParams
		FromDate         string `query:"from_date" required:"true"`
		ToDate           string `query:"to_date" required:"true"`
		Account          string `query:"account" required:"false"`
		VoucherNo        string `query:"voucher_no" required:"false"`
		AgainstVoucherNo string `query:"agains_voucher_no" required:"false"`
		PartyType        string `query:"party_type" required:"false"`
		Party            string `query:"party" required:"false"`
		Currency         string `query:"currency" required:"false"`
		Project          string `query:"project" required:"false"`
		CostCenter       string `query:"cost_center" required:"false"`
	}

	RequestAccountPayable struct {
		AuthParams
		FromDate     string `query:"from_date" required:"true"`
		ToDate       string `query:"to_date" required:"true"`
		Party        string `query:"party" required:"false"`
		ProjectID    string `query:"project_id" required:"false"`
		CostCenterID string `query:"cost_center_id" required:"false"`
	}
	RequestAccountReceivable struct {
		AuthParams
		FromDate     string `query:"from_date" required:"true"`
		ToDate       string `query:"to_date" required:"true"`
		Party        string `query:"party" required:"false"`
		ProjectID    string `query:"project_id" required:"false"`
		CostCenterID string `query:"cost_center_id" required:"false"`
	}

	SumaryEntryDto struct {
		PartyType           string `json:"party_type"`
		PartyName           string `json:"party_name"`
		PartyUUID           string `json:"party_uuid"`
		TotalInvoicedAmount int    `json:"total_invoiced_amount"`
		TotalPaidAmount     int    `json:"total_paid_amount"`
		Currency            string `json:"currency"`
	}

	AccountPayableEntryDto struct {
		PostingDate           time.Time `json:"posting_date"`
		PartyType             string    `json:"party_type"`
		PartyName             string    `json:"party_name"`
		PartyUUID             string    `json:"party_uuid"`
		ReceivableAccount     string    `json:"receivable_account"`
		ReceivableAccountUUID string    `json:"receivable_account_uuid"`
		VoucherType           string    `json:"voucher_type"`
		VoucherNo             string    `json:"voucher_no"`
		DueDate               string    `json:"due_date"`
		InvoicedAmount        int       `json:"invoiced_amount"`
		PaidAmount            int       `json:"paid_amount"`
		Currency              string    `json:"currency"`
	}

	AccountReceivableEntryDto struct {
		PostingDate           time.Time `json:"posting_date"`
		PartyType             string    `json:"party_type"`
		PartyName             string    `json:"party_name"`
		PartyUUID             string    `json:"party_uuid"`
		ReceivableAccount     string    `json:"receivable_account"`
		ReceivableAccountUUID string    `json:"receivable_account_uuid"`
		VoucherType           string    `json:"voucher_type"`
		VoucherNo             string    `json:"voucher_no"`
		DueDate               string    `json:"due_date"`
		InvoicedAmount        int       `json:"invoiced_amount"`
		PaidAmount            int       `json:"paid_amount"`
		Currency              string    `json:"currency"`
	}
	GeneralLedgerData struct {
		Entries []GeneralLedgerEntryDto `json:"entries"`
		Opening GeneralLedgerOpening    `json:"opening"`
	}

	GeneralLedgerOpening struct {
		Credit         int `json:"credit"`
		Debit          int `json:"debit"`
		OpeningBalance int `json:"opening_balance"`
	}

	GeneralLedgerEntryDto struct {
		PostingDate    time.Time `json:"posting_date"`
		Account        string    `json:"account"`
		Debit          int       `json:"debit"`
		Credit         int       `json:"credit"`
		Balance        int       `json:"balance"`
		AgainstAccount string    `json:"against_account"`
		VoucherType    string    `json:"voucher_type"`
		VoucherSubtype string    `json:"voucher_subtype"`
		VoucherNo      string    `json:"voucher_no"`
		PartyType      string    `json:"party_type"`
		PartyName      string    `json:"party_name"`
		Currency       string    `json:"currency"`
	}
)
