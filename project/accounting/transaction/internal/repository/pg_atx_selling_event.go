package transaction_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"time"
)

type AtxSellingEventRepo interface {
	OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) error
	OnDeleveryNoteSubmitted(ctx context.Context, payload event.StatusReceiptEventData) error
}

type atxSellingEventRepo struct {
	transactioRepo TransactionRepository
	core           repository.CoreService
	accounting     repository.AccountingService
	currency       helpers.CurrencyHelper
}

func NewAtxSellingEventRepo(
	transactioRepo TransactionRepository,
	core repository.CoreService,
	accounting repository.AccountingService,
	helpers *helpers.Helpers,
) AtxSellingEventRepo {
	return &atxSellingEventRepo{
		transactioRepo: transactioRepo,
		core:           core,
		accounting:     accounting,
		currency:       helpers.Currency,
	}
}

func (r *atxSellingEventRepo) OnDeleveryNoteSubmitted(ctx context.Context, payload event.StatusReceiptEventData) (err error) {
	tx := payload.Tx
	receipt := payload.Receipt
	companyDefault := payload.CompanyDefault
	accountSetting, err := r.transactioRepo.GetAccountSettings(ctx, tx, payload.Receipt.CompanyID)
	if err != nil {
		return err
	}
	stockSetting, err := r.transactioRepo.GetStockSettings(ctx, tx, payload.Receipt.CompanyID)
	if err != nil {
		return err
	}
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefault, receipt.Currency, true, false)
	if err != nil {
		return
	}

	//Get total Amount Receipt
	totalAmount, err := r.getAmountFromLineItems(ctx, tx, receipt.Code)
	if err != nil {
		return err
	}

	totalAmount = r.currency.CurrencyExchange(int64(totalAmount), exchangeRate)

	transactions := r.transactionInventory(accountSetting, stockSetting, totalAmount,
		payload.ReceiptPartyType, receipt.Code, receipt.CostCenterID, receipt.ProjectID, receipt.PostingTime,
		receipt.PostingDate,companyDefault.Currency)

	err = r.transactioRepo.SaveTransactionsBatch(ctx, tx, transactions)
	return
}

func (r *atxSellingEventRepo) transactionInventory(accountSetting model.AccountSetting, stockSetting model.StockSetting,
	lineItemAmount int64, voucherType, voucherCode string, costCenterID, projectID *int64,
	postingTime string, postingDate time.Time,currency string,
) []*model.TransactionLedger {
	credit := model.TransactionLedger{
		Ledger:        stockSetting.InventoryAccount,
		LedgerAgainst: &accountSetting.CostOfGoodSoldAccount,
		Credit:        lineItemAmount,
		VoucherCode:   voucherCode,
		VoucherType:   voucherType,
		ProjectID:     projectID,
		CostCenterID:  costCenterID,
		PostingTime:   postingTime,
		PostingDate:   postingDate,
		Currency: currency,
	}
	debit := model.TransactionLedger{
		Ledger:        accountSetting.CostOfGoodSoldAccount,
		LedgerAgainst: &stockSetting.InventoryAccount,
		Debit:         lineItemAmount,
		VoucherCode:   voucherCode,
		VoucherType:   voucherType,
		ProjectID:     projectID,
		CostCenterID:  costCenterID,
		PostingTime:   postingTime,
		PostingDate:   postingDate,
		Currency: currency,
	}
	return []*model.TransactionLedger{&credit, &debit}
}

func (r *atxSellingEventRepo) getAmountFromLineItems(ctx context.Context, tx *query.QueryTx,
	voucherCode string) (totalAmount int64, err error) {
	var serialNos []model.SerialNo
	snTxQ := tx.SerialNoTransaction
	serialNoQ := tx.SerialNo
	batchBundleQ := tx.BatchBundle
	err = snTxQ.WithContext(ctx).
		Select(serialNoQ.ValuationRate).
		Join(serialNoQ, snTxQ.SerialNoID.EqCol(serialNoQ.ID)).
		Join(batchBundleQ, snTxQ.BatchBundleID.EqCol(batchBundleQ.ID)).
		Where(
			batchBundleQ.VoucherCode.Eq(voucherCode),
		).
		Order(snTxQ.ID.Asc()).
		Scan(&serialNos)
	for _, serialNo := range serialNos {
		totalAmount += int64(serialNo.ValuationRate)
	}
	return
}

func (r *atxSellingEventRepo) OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) (err error) {
	tx := payload.Tx
	invoice := payload.Invoice
	currency := payload.CompanyDefaults.Currency
	accountSetting, err := r.transactioRepo.GetAccountSettings(ctx, tx, payload.Invoice.CompanyID)
	if err != nil {
		return
	}
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, payload.CompanyDefaults, invoice.Currency, true, false)
	if err != nil {
		return
	}
	totalItemAmount := r.currency.CurrencyExchange(payload.LineItemsData.TotalAmount, exchangeRate)
	tacTransactions, totalTaxAmount, err := r.transactioRepo.ProcessTaxLines(ctx, tx, invoice.ID, invoice.PostingDate,
		invoice.PostingTime,currency, invoice.Code, proto.PartyType_saleInvoice.String(),
		proto.VoucherSubtype_debitNote.String(), invoice.CostCenterID, invoice.ProjectID, payload.TaxLinesData.TaxLines,
	exchangeRate)
	if err != nil {
		return
	}
	totalAmount := totalItemAmount + totalTaxAmount
	var transactions []*model.TransactionLedger
	receivableEntry := model.TransactionLedger{
		Ledger:         accountSetting.ReceivableAccount,
		LedgerAgainst:  &accountSetting.IncomeAccount,
		Debit:          totalAmount,
		VoucherCode:    invoice.Code,
		PartyID:        &invoice.PartyID,
		VoucherType:    payload.InvoicePartyType,
		VoucherSubtype: proto.VoucherSubtype_debitNote.String(),
		ProjectID:      invoice.ProjectID,
		CostCenterID:   invoice.CostCenterID,
		PostingDate: invoice.PostingDate,
		PostingTime: invoice.PostingTime,
		Currency: currency,
	}
	incomeEntry := model.TransactionLedger{
		Ledger:         accountSetting.IncomeAccount,
		LedgerAgainst:  &accountSetting.ReceivableAccount,
		Credit:         totalItemAmount,
		VoucherCode:    invoice.Code,
		VoucherType:    payload.InvoicePartyType,
		VoucherSubtype: proto.VoucherSubtype_debitNote.String(),
		ProjectID:      invoice.ProjectID,
		CostCenterID:   invoice.CostCenterID,
		PostingDate: invoice.PostingDate,
		PostingTime: invoice.PostingTime,
		Currency: currency,
	}
	if invoice.UpdateStock {
		//Get total Amount Receipt
		stockAmount, err := r.getAmountFromLineItems(ctx, tx, invoice.Code)
		if err != nil {
			return err
		}
		stockAmount = r.currency.CurrencyExchange(int64(stockAmount), exchangeRate)
		if stockAmount > 0 {
			stockSetting, err := r.transactioRepo.GetStockSettings(ctx, tx, payload.Invoice.CompanyID)
			if err != nil {
				return err
			}
			transactions = r.transactionInventory(accountSetting, stockSetting, stockAmount,
				payload.InvoicePartyType, invoice.Code, invoice.CostCenterID, invoice.ProjectID,
				invoice.PostingTime, invoice.PostingDate,currency)
		}
	}

	transactions = append(transactions, &receivableEntry, &incomeEntry)
	transactions = append(transactions, tacTransactions...)
	err = r.transactioRepo.SaveTransactionsBatch(ctx, tx, transactions)
	return
}

type DeliveryNoteItemLine struct {
	Quantity          int32
	ItemID            int64
	AcceptedWarehouse int64
}
