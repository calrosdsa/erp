package dto

import "time"

type (
	CreateTransaction struct {
		LedgerID   int64    `json:"ledger_id"`
		LedgerIDDr int64    `json:"ledger_id_dr"`
		Amount     int32     `json:"amount"`
		DateTime   time.Time `json:"datetime"`

		VoucherCode    string `json:"voucher_code"`
		VoucherType    string `json:"voucher_type"`
		VoucherSubtype string `json:"voucher_subtype"`
		PartyID        *int64 `json:"party_id"`
	}
)
