package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CashOutflowsRequest struct {
		DefaultListParams

		CreatedAt string `query:"created_at" required:"false"`
		UpdatedAt string `query:"updated_at" required:"false"`
	}
	CashOutflowDataRequest struct {
		Body CashOutflowData
	}
	CashOutflowData struct {
		ID                  int64               `json:"id"`
		Fields              CashOutflowFields   `json:"fields"`
		CreateTaxAndCharges CreateTaxAndChanges `json:"tax_and_charges"`
	}

	CashOutflowFields struct {
		PartyID         int64      `json:"party_id"`
		PartyType       string     `json:"party_type"`
		Concept         *string    `json:"concept" required:"false"`
		CashOutflowType *string    `json:"cash_outflow_type" required:"false"`
		InvoiceNo       *string    `json:"invoice_no" required:"false"`
		Nit             *string    `json:"nit" required:"false"`
		AuthCode        *string    `json:"auth_code" required:"false"`
		CtrlCode        *string    `json:"ctrl_code" required:"false"`
		EmisionDate     *time.Time `json:"emision_date" required:"false"`
		PostingDate     time.Time  `json:"posting_date"`
		PostingTime     string     `json:"posting_time"`
		Tz              string     `json:"tz"`
		Amount          int64      `json:"amount"`
		ProjectID       *int64     `json:"project_id" required:"false"`
		CostCenterID    *int64     `json:"cost_center_id" required:"false"`
	}

	CashOutflowDto struct {
		ID     int64  `json:"id"`
		Status string `json:"status"`
		Code   string `json:"code"`

		PartyID   int64  `json:"party_id"`
		Party     string `json:"party"`
		PartyUUID string `json:"party_uuid"`

		PartyType       string     `json:"party_type"`
		Concept         *string    `json:"concept"`
		CashOutflowType *string    `json:"cash_outflow_type"`
		InvoiceNo       *string    `json:"invoice_no"`
		Nit             *string    `json:"nit"`
		AuthCode        *string    `json:"auth_code"`
		CtrlCode        *string    `json:"ctrl_code"`
		EmisionDate     *time.Time `json:"emision_date"`
		PostingDate     time.Time  `json:"posting_date"`
		PostingTime     string     `json:"posting_time"`
		Tz              string     `json:"tz"`
		Amount          int64      `json:"amount"`

		AccountingDimensionDto
	}
)

func CashOutFlowModel(m model.CashOutflow) CashOutflowDto {
	r := CashOutflowDto{}
	r.ID = m.ID
	r.Code = m.Code
	return r
}
