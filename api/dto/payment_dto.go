package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	RequestPayments struct {
		PaginationParams
		AuthParams
		OptionalQueryParams

		InvoiceID         string `query:"invoice_id" required:"false"`
		PostingDate       string `query:"posting_date" required:"false"`
		PaymentType       string `query:"payment_type" required:"false"`
		Amount            string `query:"amount" required:"false"`
		Code              string `query:"code" required:"false"`
		PartyID           string `query:"party_id" required:"false"`
		AccountPaidFromID string `query:"account_paid_from_id" required:"false"`
		AccountPaidToID   string `query:"account_paid_to_id" required:"false"`
	}
	CreatePaymentRequest struct {
		Body PaymentBody
	}
	EditPaymentRequest struct {
		Body PaymentBody
	}

	PaymentBody struct {
		PaymentData         PaymentData              `json:"payment"`
		CreateTaxAndCharges CreateTaxAndChanges      `json:"tax_and_charges"`
		PaymentReferences   []CreatePaymentReference `json:"payment_references"`
	}

	PaymentData struct {
		ID     int64         `json:"id" required:"false"`
		Fields PaymentFields `json:"fields" required:"true"`
		// ModeOfPayment string `json:"mode_of_payment" required:"true"`
	}

	PaymentFields struct {
		PostingDate time.Time `json:"posting_date" required:"true"`
		PaymentType string    `json:"payment_type" required:"true"`
		Amount      int64     `json:"amount" required:"true"`

		AccountPaidFromID int64 `json:"account_paid_from_id" required:"true"`
		AccountPaidToID   int64 `json:"account_paid_to_id" required:"true"`

		PartyID              int64  `json:"party_id"`
		PartyBankAccountID   *int64 `json:"party_bank_account_id" required:"false"`
		CompanyBankAccountID *int64 `json:"company_bank_account_id" required:"false"`

		ProjectID    *int64 `json:"project_id" required:"false"`
		CostCenterID *int64 `json:"cost_center_id" required:"false"`

		ChequeReferenceNo *string `json:"cheque_reference_no" required:"false"`
		ChequeReferenceDate *time.Time  `json:"cheque_reference_date" required:"false"`
	}

	CreatePaymentReference struct {
		PartyID     int64   `json:"party_id"`
		PartyCode   string  `json:"party_code"`
		PartyType   string  `json:"party_type"`
		Total       float64 `json:"total"`
		Outstanding float64 `json:"outstanding"`
		Allocated   float64 `json:"allocated"`
		Currency    string  `json:"currency"`
	}

	PaymentDto struct {
		ID          int64     `json:"id"`
		Code        string    `json:"code"`
		Amount      int64       `json:"amount"`
		PostingDate time.Time `json:"posting_date"`
		CreatedAt   time.Time `json:"created_at"`
		PaymentType string    `json:"payment_type"`
		Status      string    `json:"status"`

		//Party
		PartyID   int64  `json:"party_id"`
		PartyName string `json:"party_name"`
		PartyUUID string `json:"party_uuid"`
		PartyType string `json:"party_type"`
		//Company Bank Account
		CompanyBankAccount     *string `json:"company_bank_account"`
		CompanyBankAccountID   *int64  `json:"company_bank_account_id"`
		CompanyBankAccountUUID *string `json:"company_bank_account_uuid"`

		//Party Bank Account
		PartyBankAccount     *string `json:"party_bank_account"`
		PartyBankAccountID   *int64  `json:"party_bank_account_id"`
		PartyBankAccountUUID *string `json:"party_bank_account_uuid"`

		ChequeReferenceNo *string `json:"cheque_reference_no"`
		ChequeReferenceDate *time.Time `json:"cheque_reference_date"`

		AccountingDimensionDto
	}

	PaymentReferenceDto struct {
		PartyID     int64  `json:"party_id"`
		PartyCode   string `json:"party_code"`
		PartyType   string `json:"party_type"`
		Total       int64  `json:"total"`
		Outstanding int64  `json:"outstanding"`
		Allocated   int64  `json:"allocated"`
		Currency    string `json:"currency"`
	}

	PaymentDetailDto struct {
		PaymentDto

		PaidFromID       int64  `json:"paid_from_id"`
		PaidFromName     string `json:"paid_from_name"`
		PaidFromUUID     string `json:"paid_from_uuid"`
		PaidFromCurrency string `json:"paid_from_currency"`

		PaidToID       int64  `json:"paid_to_id"`
		PaidToName     string `json:"paid_to_name"`
		PaidToUUID     string `json:"paid_to_uuid"`
		PaidToCurrency string `json:"paid_to_currency"`

		PaymentReferences []PaymentReferenceDto `json:"payment_references" gorm:"-"`
	}

	PaymentAccountsDto struct {
		ReceivableAcctID       int64  `json:"receivable_acct_id"`
		ReceivableAcct         string `json:"receivable_acct"`
		ReceivableAcctCurrency string `json:"receivable_acct_currency"`

		PayableAcctID       int64  `json:"payable_acct_id"`
		PayableAcct         string `json:"payable_acct"`
		PayableAcctCurrency string `json:"payable_acct_currency"`

		CashAcctID       int64  `json:"cash_acct_id"`
		CashAcct         string `json:"cash_acct"`
		CashAcctCurrency string `json:"cash_acct_currency"`
	}
)

func PaymentDtoFromModel(m *model.Payment) PaymentDto {
	r := PaymentDto{}
	r.Amount = m.Amount
	r.Code = m.Code
	r.CreatedAt = m.CreatedAt
	r.PaymentType = m.PaymentType
	r.PostingDate = m.PostingDate
	r.Status = m.Status
	return r
}
