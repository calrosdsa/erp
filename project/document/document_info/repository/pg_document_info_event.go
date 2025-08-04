package documentinfo_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain/event"
)

type DocumentInfoEventRepo interface {
	CreateDocumentInfoForOrder(ctx context.Context, payload event.OrderEventData) (err error)
	CreateDocumentInfoForInvoice(ctx context.Context, payload event.InvoiceEventData) (err error)
	CreateDocumentInfoForReceipt(ctx context.Context, payload event.ReceiptEventData) (err error)
	CreateDocumentInfoForCashOutflow(ctx context.Context, payload event.CashOutflowEventData) (err error)
}

type documentInfoRepo struct {
}

func NewDocumentEventRepository() DocumentInfoEventRepo {
	return &documentInfoRepo{}
}

func (r *documentInfoRepo) createDocTerms(ctx context.Context, tx *query.QueryTx, docID int64) (err error) {
	docTermsQ := tx.DocTerm
	docTerm := &model.DocTerm{}
	docTerm.DocID = docID
	err = docTermsQ.WithContext(ctx).Save(docTerm)
	return
}
func (r *documentInfoRepo) createDocAccounts(ctx context.Context, tx *query.QueryTx, docID int64,
	docPartyType string, companyID int64) (err error) {
	docAccountQ := tx.DocAccount
	docAccount := &model.DocAccount{}
	docAccount.DocID = docID
	accountDefaults, err := tx.AccountSetting.Where(
		tx.AccountSetting.CompanyID.Eq(companyID),
	).First()
	if err != nil {
		return
	}
	stockSetting, err := tx.StockSetting.Where(
		tx.StockSetting.CompanyID.Eq(companyID),
	).First()
	if err != nil {
		return
	}
	switch docPartyType {
	case proto.PartyType_purchaseInvoice.String():
		docAccount.CreditAccountID = &accountDefaults.PayableAccount
		// if updateStock {
		// 	docAccount.DebitAccountID = &stockSetting.InventoryAccount
		// }else {
		// 	docAccount.DebitAccountID = &stockSetting.StockReceivedButNotBilled
		// }
	case proto.PartyType_saleInvoice.String():
		docAccount.DebitAccountID = &accountDefaults.ReceivableAccount
		
	case proto.PartyType_cashOutflow.String():
		docAccount.CreditAccountID = &accountDefaults.CashAccunt
		docAccount.DebitAccountID = &accountDefaults.PayableAccount
	case proto.PartyType_deliveryNote.String():
		docAccount.CreditAccountID = &stockSetting.InventoryAccount
		docAccount.DebitAccountID = &accountDefaults.CostOfGoodSoldAccount
	case proto.PartyType_purchaseReceipt.String():
		docAccount.CreditAccountID = &stockSetting.InventoryAccount
		docAccount.DebitAccountID = &stockSetting.StockReceivedButNotBilled
	}
	err = docAccountQ.WithContext(ctx).Save(docAccount)
	return
}

func (r *documentInfoRepo) createDoAddressAndContact(ctx context.Context, tx *query.QueryTx, docID int64) (err error) {
	aacQ := tx.AddressAndContact
	addressAndContact := &model.AddressAndContact{}
	addressAndContact.DocID = docID
	err = aacQ.WithContext(ctx).Save(addressAndContact)
	return
}

func (r *documentInfoRepo) CreateDocumentInfoForCashOutflow(ctx context.Context, payload event.CashOutflowEventData) (
	err error) {
	tx := payload.Tx
	if err = r.createDoAddressAndContact(ctx, tx, payload.CashOutflow.ID); err != nil {
		return
	}
	if err = r.createDocTerms(ctx, tx, payload.CashOutflow.ID); err != nil {
		return
	}

	if err = r.createDocAccounts(ctx, tx, payload.CashOutflow.ID, proto.PartyType_cashOutflow.String(),
		payload.CashOutflow.CompanyID); err != nil {
		return
	}
	return
}

func (r *documentInfoRepo) CreateDocumentInfoForOrder(ctx context.Context, payload event.OrderEventData) (
	err error) {
	tx := payload.Tx
	if err = r.createDoAddressAndContact(ctx, tx, payload.Order.ID); err != nil {
		return
	}
	if err = r.createDocTerms(ctx, tx, payload.Order.ID); err != nil {
		return
	}
	return
}

func (r *documentInfoRepo) CreateDocumentInfoForInvoice(ctx context.Context, payload event.InvoiceEventData) (
	err error) {
	tx := payload.Tx
	if err = r.createDoAddressAndContact(ctx, tx, payload.Invoice.ID); err != nil {
		return
	}
	if err = r.createDocTerms(ctx, tx, payload.Invoice.ID); err != nil {
		return
	}

	if err = r.createDocAccounts(ctx, tx, payload.Invoice.ID, payload.InvoicePartyType, payload.Invoice.CompanyID); err != nil {
		return
	}
	return
}

func (r *documentInfoRepo) CreateDocumentInfoForReceipt(ctx context.Context, payload event.ReceiptEventData) (
	err error) {
	tx := payload.Tx
	if err = r.createDoAddressAndContact(ctx, tx, payload.Receipt.ID); err != nil {
		return
	}
	if err = r.createDocTerms(ctx, tx, payload.Receipt.ID); err != nil {
		return
	}
	if err = r.createDocAccounts(ctx, tx, payload.Receipt.ID, payload.ReceiptPartyType, payload.Receipt.CompanyID); err != nil {
		return
	}
	return
}
