package tac_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"
)

type TacRepository interface {
	GetTACLines(req *common.RequestContext, d *dto.RequestTaxLines) (
		[]dto.TaxAndChargeLineDto, error)
	EditTaxAndChargeLine(req *common.RequestContext,
		d *dto.EditTaxLineRequest) (err error)
	CreateTaxAndChargeLine(req *common.RequestContext,
		d *dto.AddTaxLineRequest) (err error)
	DeleteTaxAndChargeLine(req *common.RequestContext, d *dto.DeleteTaxLineRequest) (err error)
}

type tacRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	currency  helpers.CurrencyHelper
}

func NewTacRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) TacRepository {
	return &tacRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		currency:  helpers.Currency,
	}
}
func (r *tacRepository) DeleteTaxAndChargeLine(req *common.RequestContext, d *dto.DeleteTaxLineRequest) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = r.Q.TaxAndChargeLine.WithContext(req.Ctx).Where(
		r.Q.TaxAndChargeLine.ID.Eq(d.Body.ID),
	).Delete()
	err = r.updateDocTotalAmount(req, tx, d.Body.DocPartyID, d.Body.DocPartyType, d.Body.TotalAmount)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *tacRepository) GetTACLines(req *common.RequestContext, d *dto.RequestTaxLines) (
	[]dto.TaxAndChargeLineDto, error) {
	docID := r.convertor.StrtoInt(d.ID)
	tacLineQ := r.Q.TaxAndChargeLine
	ledgerQ := r.Q.Ledger
	tacLines := []dto.TaxAndChargeLineDto{}
	err := tacLineQ.WithContext(req.Ctx).Select(
		tacLineQ.ID, tacLineQ.DocPartyID, tacLineQ.Type, tacLineQ.TaxRate,
		tacLineQ.Amount, tacLineQ.IsDeducted,
		ledgerQ.ID.As("account_head_id"), ledgerQ.Name.As("account_head"), ledgerQ.UUID.As("account_head_uuid"),
	).
		Join(ledgerQ, ledgerQ.ID.EqCol(tacLineQ.AccountHead)).
		Where(tacLineQ.DocPartyID.Eq(docID)).
		Order(tacLineQ.ID.Asc()).
		Scan(&tacLines)
	return tacLines, err
}

func (r *tacRepository) EditTaxAndChargeLine(req *common.RequestContext,
	d *dto.EditTaxLineRequest) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	lineData := d.Body
	tacLineQ := r.Q.TaxAndChargeLine
	_, err = r.Q.TaxAndChargeLine.WithContext(req.Ctx).
		Where(tacLineQ.ID.Eq(d.Body.ID)).UpdateSimple(
		tacLineQ.AccountHead.Value(lineData.LedgerID),
		tacLineQ.TaxRate.Value(lineData.TaxRate),
		tacLineQ.Amount.Value(r.currency.FloatToInt(lineData.Amount)),
		tacLineQ.Type.Value(lineData.Type),
		tacLineQ.IsDeducted.Value(lineData.IsDeducted),
	)
	err = r.updateDocTotalAmount(req, tx, d.Body.DocPartyID, d.Body.DocPartyType, d.Body.TotalAmount)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func(r *tacRepository) adjustTaxAndCharges(ctx context.Context,tx *query.QueryTx,docPartyID int64)(err error){
	tacLineQ := tx.TaxAndChargeLine
	_,err =tacLineQ.WithContext(ctx).Where(tacLineQ.DocPartyID.Eq(docPartyID)).Find()
	if err != nil {
		return 
	}
	// for i,tacLines:= ta
	return
}

func (r *tacRepository) CreateTaxAndChargeLine(req *common.RequestContext,
	d *dto.AddTaxLineRequest) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	fmt.Println("Create", d.Body)
	taxAndCharge := model.TaxAndChargeLine{}
	taxAndCharge.AccountHead = d.Body.LedgerID
	taxAndCharge.DocPartyID = d.Body.DocPartyID
	taxAndCharge.Amount = r.currency.FloatToInt(d.Body.Amount)
	taxAndCharge.IsDeducted = d.Body.IsDeducted
	if d.Body.TaxRate != 0 {
		taxAndCharge.TaxRate = &d.Body.TaxRate
	}
	taxAndCharge.Type = d.Body.Type
	err = tx.TaxAndChargeLine.WithContext(req.Ctx).Save(&taxAndCharge)
	if err != nil {
		return
	}
	err = r.updateDocTotalAmount(req, tx, d.Body.DocPartyID, d.Body.DocPartyType, d.Body.TotalAmount)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *tacRepository) updateDocTotalAmount(req *common.RequestContext,
	tx *query.QueryTx, docPartyID int64, docPartyType string, totalAmount float64) (err error) {
	switch docPartyType {
	case proto.PartyType_purchaseOrder.String(), proto.PartyType_saleOrder.String():
		fmt.Println("UPDATE DOC TOTAL AMOUNT", docPartyID, docPartyType)
		_, err = tx.ProgressOrder.WithContext(req.Ctx).Where(
			tx.ProgressOrder.OrderID.Eq(docPartyID),
		).UpdateSimple(
			tx.ProgressOrder.TotalAmount.Value(r.currency.FloatToInt(totalAmount)),
		)
	case proto.PartyType_purchaseInvoice.String(), proto.PartyType_saleInvoice.String():
		_, err = tx.ProgressInvoice.WithContext(req.Ctx).Where(
			tx.ProgressInvoice.InvoiceID.Eq(docPartyID),
		).UpdateSimple(
			tx.ProgressInvoice.TotalAmount.Value(r.currency.FloatToInt(totalAmount)),
		)
	}
	return err
}
