package tac_repo

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"fmt"
)

type TacRepositoryEvent interface {
	OnQuotationCreated(ctx context.Context, payload event.QuotationEventData) (err error)
	OnOrderCreated(ctx context.Context, payload event.OrderEventData) (err error)
	OnInvoiceCreated(ctx context.Context, payload event.InvoiceEventData) (err error)
	OnReceiptCreated(ctx context.Context, payload event.ReceiptEventData) (err error)
	OnChargesTemplateCreated(ctx context.Context, payload event.ChargesTemplateEventData) (err error)
	OnPaymentCreated(ctx context.Context, payload event.PaymentEventData) (err error)
	EditTaxAndChargeLines(tx *query.QueryTx, ctx context.Context,
		charges dto.CreateTaxAndChanges, docPartyID int64) error
	CreateTaxAndChargeLines(tx *query.QueryTx, ctx context.Context,
		d dto.CreateTaxAndChanges, docPartyID int64) (err error)
		
}

type tacRepositoryEvent struct {
	currency helpers.CurrencyHelper
}

func NewTacEventRepository(
	helpers *helpers.Helpers,
) TacRepositoryEvent {
	return &tacRepositoryEvent{
		currency: helpers.Currency,
	}
}

func (r *tacRepositoryEvent) OnPaymentCreated(ctx context.Context, payload event.PaymentEventData) (err error) {
	tx := payload.Tx
	err = r.CreateTaxAndChargeLines(tx, ctx, payload.Body.CreateTaxAndCharges, payload.Payment.ID)
	return
}
func (r *tacRepositoryEvent) OnOrderCreated(ctx context.Context, payload event.OrderEventData) (
	err error) {
	tx := payload.Tx
	err = r.CreateTaxAndChargeLines(tx, ctx, payload.Body.CreateTaxAndCharges, payload.Order.ID)
	return
}
func (r *tacRepositoryEvent) OnInvoiceCreated(ctx context.Context, payload event.InvoiceEventData) (
	err error) {
	tx := payload.Tx
	err = r.CreateTaxAndChargeLines(tx, ctx, payload.Body.CreateTaxAndCharges, payload.Invoice.ID)
	return
}
func (r *tacRepositoryEvent) OnReceiptCreated(ctx context.Context, payload event.ReceiptEventData) (
	err error) {
	tx := payload.Tx
	err = r.CreateTaxAndChargeLines(tx, ctx, payload.Body.CreateTaxAndCharges, payload.Receipt.ID)
	return
}

func (r *tacRepositoryEvent) OnQuotationCreated(ctx context.Context, payload event.QuotationEventData) (
	err error) {
	tx := payload.Tx
	err = r.CreateTaxAndChargeLines(tx, ctx, payload.Body.CreateTaxAndCharges, payload.Quotation.ID)
	return
}
func (r *tacRepositoryEvent) OnChargesTemplateCreated(ctx context.Context, payload event.ChargesTemplateEventData) (
	err error) {
	tx := payload.Tx
	err = r.CreateTaxAndChargeLines(tx, ctx, payload.ChargeTemplateData.CreateTaxAndCharges, payload.ChargesTemplate.ID)
	return
}

func (r *tacRepositoryEvent) EditTaxAndChargeLines(tx *query.QueryTx, ctx context.Context,
	charges dto.CreateTaxAndChanges, docPartyID int64) error {
	_, err := tx.WithContext(ctx).TaxAndChargeLine.Unscoped().Where(
		tx.TaxAndChargeLine.DocPartyID.Eq(docPartyID),
	).Delete()
	if err != nil {
		return err
	}
	err = r.CreateTaxAndChargeLines(tx, ctx, charges, docPartyID)
	return err
}

func (r *tacRepositoryEvent) deleteTaxAndChargeLines(tx *query.QueryTx, ctx context.Context, docPartyID int64) error {
	_, err := tx.WithContext(ctx).TaxAndChargeLine.Unscoped().Where(
		tx.TaxAndChargeLine.DocPartyID.Eq(docPartyID),
	).Delete()
	if err != nil {
		return err
	}

	return err
}

func (r *tacRepositoryEvent) CreateTaxAndChargeLines(tx *query.QueryTx, ctx context.Context,
	d dto.CreateTaxAndChanges, docPartyID int64) (err error) {
	taxAndCharges := make([]*model.TaxAndChargeLine, len(d.TaxAndCharges))
	for i, line := range d.TaxAndCharges {
		taxAndCharge := &model.TaxAndChargeLine{}
		fmt.Println("TAX AND LINE IS DEDUCTED", line.IsDeducted)
		taxAndCharge.DocPartyID = docPartyID
		taxAndCharge.AccountHead = line.LedgerID
		taxAndCharge.IsDeducted = line.IsDeducted
		taxAndCharge.Amount = r.currency.FloatToInt(line.Amount)
		if line.TaxRate != 0 {
			taxAndCharge.TaxRate = &line.TaxRate
		}
		taxAndCharge.Type = line.Type
		taxAndCharges[i] = taxAndCharge
	}
	err = tx.TaxAndChargeLine.WithContext(ctx).CreateInBatches(taxAndCharges, len(taxAndCharges))
	return
}

// func (r *tacRepositoryEvent) CreateTaxAndChangeLine(tx *query.QueryTx,req *common.RequestContext,
// 	d dto.TaxAndChargeLineData,docPartyID int64)(err error){
// 	taxAndCharge := model.TaxAndChargeLine{}
// 	taxAndCharge.AccountHead = d.LedgerID
// 	taxAndCharge.DocPartyID = docPartyID
// 	if d.TaxRate != 0 {
// 		taxAndCharge.TaxRate =&d.TaxRate
// 	}
// 	taxAndCharge.Type = d.Type
// 	err =tx.TaxAndChargeLine.WithContext(req.Ctx).Save(&taxAndCharge)
// 	return
// }
