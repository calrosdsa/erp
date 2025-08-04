package dto

import "erp/gen/db/model"

type (
	RequestTaxLines struct {
		ID string `query:"id" required:"true"`
	}
	CreateTaxAndChanges struct {
		TaxAndCharges []TaxAndChargeLineData `json:"lines"`
	}
	EditTaxLineRequest struct {
		Body struct {
			ID           int32   `json:"id"`
			DocPartyID   int64   `json:"doc_party_id" required:"false"`
			DocPartyType string  `json:"doc_party_type" required:"false"`
			TotalAmount  float64 `json:"total_amount"`
			TaxAndChargeLineData
		}
	}

	AddTaxLineRequest struct {
		Body struct {
			DocPartyID   int64   `json:"doc_party_id" required:"false"`
			DocPartyType string  `json:"doc_party_type" required:"false"`
			TotalAmount  float64 `json:"total_amount"`
			TaxAndChargeLineData
		}
	}
	DeleteTaxLineRequest struct {
		Body struct {
			ID           int32   `json:"id"`
			DocPartyID   int64   `json:"doc_party_id" required:"false"`
			DocPartyType string  `json:"doc_party_type" required:"false"`
			TotalAmount  float64 `json:"total_amount"`
		}
	}

	TaxAndChargeLineData struct {
		Type       string  `json:"type" required:"true"`
		LedgerID   int64   `json:"ledger" required:"true"`
		TaxRate    int16   `json:"tax_rate" requied:"true"`
		Amount     float64 `json:"amount" required:"true"`
		IsDeducted bool    `json:"is_deducted" required:"true"`
	}
	TaxAndChargeLineDto struct {
		ID              int64  `json:"id"`
		TaxRate         int16  `json:"tax_rate"`
		Amount          int32  `json:"amount"`
		IsDeducted      bool   `json:"is_deducted"`
		Type            string `json:"type"`
		AccountHeadID   int    `json:"account_head_id"`
		AccountHead     string `json:"account_head"`
		AccountHeadUUID string `json:"account_head_uuid"`
	}
	TaxLinesData  struct {
		TaxLines []TaxLineWhitAccountHead
		TotalAmount int64
	}

	TaxLineWhitAccountHead struct {
		model.TaxAndChargeLine
		AcctHeadRootType string
		IsOffsetAccount  bool
	}
)
