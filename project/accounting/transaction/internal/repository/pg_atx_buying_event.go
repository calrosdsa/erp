package transaction_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"fmt"
	"time"
)

type AtxBuyingEventRepo interface {
	OnReceiptSubmitted(ctx context.Context, payload event.StatusReceiptEventData) error
	OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) error
}

type atxBuyingEventRepo struct {
	transactioRepo TransactionRepository
	core           repository.CoreService
	accounting     repository.AccountingService
	currency       helpers.CurrencyHelper
}

func NewAtxBuyingRepo(
	transactioRepo TransactionRepository,
	core repository.CoreService,
	accounting repository.AccountingService,
	helpers *helpers.Helpers,
) AtxBuyingEventRepo {
	return &atxBuyingEventRepo{
		transactioRepo: transactioRepo,
		core:           core,
		accounting:     accounting,
		currency:       helpers.Currency,
	}
}

func (r *atxBuyingEventRepo) OnInvoiceSubmitted(ctx context.Context,
	payload event.StatusInvoiceEventData) (err error) {
	tx := payload.Tx
	invoice := payload.Invoice
	companyDefault := payload.CompanyDefaults
	currency := companyDefault.Currency
	stockSetting, err := r.transactioRepo.GetStockSettings(ctx, tx, payload.Invoice.CompanyID)
	if err != nil {
		return
	}
	accountSetting, err := r.transactioRepo.GetAccountSettings(ctx, tx, payload.Invoice.CompanyID)
	if err != nil {
		return
	}
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefault, invoice.Currency, false, true)
	if err != nil {
		return
	}

	tacTransactions, totalTaxAmount, err := r.transactioRepo.ProcessTaxLines(ctx, tx, invoice.ID, invoice.PostingDate, invoice.PostingTime,
		currency, invoice.Code, proto.PartyType_purchaseInvoice.String(),proto.VoucherSubtype_debitNote.String(),
		invoice.CostCenterID, invoice.ProjectID, payload.TaxLinesData.TaxLines,exchangeRate)
	if err != nil {
		return
	}
	var transactions []*model.TransactionLedger
	stockAmount := r.currency.CurrencyExchange(payload.LineItemsData.TotalAmount, exchangeRate)
	totalAmount := stockAmount + totalTaxAmount
	if invoice.UpdateStock {
		//Get total Amount Receipt
		payableEntry := model.TransactionLedger{
			Ledger:         accountSetting.PayableAccount,
			Credit:         totalAmount,
			VoucherCode:    invoice.Code,
			PartyID:        &invoice.PartyID,
			VoucherType:    payload.InvoicePartyType,
			VoucherSubtype: proto.VoucherSubtype_creditNote.String(),
			CostCenterID:   invoice.CostCenterID,
			PostingTime: invoice.PostingTime,
			PostingDate: invoice.PostingDate,
			ProjectID:      invoice.ProjectID,
			Currency: currency,
		}
		debitInventory := model.TransactionLedger{
			Ledger:         stockSetting.InventoryAccount,
			Debit:          stockAmount,
			VoucherCode:    invoice.Code,
			VoucherType:    payload.InvoicePartyType,
			VoucherSubtype: payload.InvoicePartyType,
			ProjectID:      invoice.ProjectID,
			CostCenterID:   invoice.CostCenterID,
			PostingTime: invoice.PostingTime,
			PostingDate: invoice.PostingDate,
			Currency: currency,
		}
		transactions = append(transactions, &payableEntry, &debitInventory)
	} else {
		payableEntry := model.TransactionLedger{
			Ledger: accountSetting.PayableAccount,
			Credit:         totalAmount,
			VoucherCode:    invoice.Code,
			PartyID:        &invoice.PartyID,
			VoucherType:    payload.InvoicePartyType,
			VoucherSubtype: proto.VoucherSubtype_creditNote.String(),
			CostCenterID:   invoice.CostCenterID,
			ProjectID:      invoice.ProjectID,
			PostingTime: invoice.PostingTime,
			PostingDate: invoice.PostingDate,
			Currency: currency,
		}
		receivedButNtBilledEntry := model.TransactionLedger{
			Ledger: stockSetting.StockReceivedButNotBilled,
			Debit:          stockAmount,
			VoucherCode:    invoice.Code,
			VoucherType:    payload.InvoicePartyType,
			VoucherSubtype: proto.VoucherSubtype_creditNote.String(),
			CostCenterID:   invoice.CostCenterID,
			ProjectID:      invoice.ProjectID,
			PostingTime: invoice.PostingTime,
			PostingDate: invoice.PostingDate,
			Currency: currency,
		}
		transactions = append(transactions, &payableEntry, &receivedButNtBilledEntry)
	}
	transactions = append(transactions, tacTransactions...)
	err = r.transactioRepo.SaveTransactionsBatch(ctx, tx, transactions)
	return err
}

func (r *atxBuyingEventRepo) OnReceiptSubmitted(ctx context.Context, payload event.StatusReceiptEventData) error {
	tx := payload.Tx
	receipt := payload.Receipt

	stockSetting, err := r.transactioRepo.GetStockSettings(ctx, tx, payload.Receipt.CompanyID)
	if err != nil {
		return err
	}

	//Get total Amount Receipt
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, payload.CompanyDefault,
		receipt.Currency, false, true)
	if err != nil {
		return err
	}

	totalAmount := r.currency.CurrencyExchange(payload.LineItemsData.TotalAmount, exchangeRate)
	fmt.Println("TOTAL AMOUNT", totalAmount)
	fmt.Println("stock setting", stockSetting)

	//Adding receipt account entry
	transactions := r.transactionInventory(stockSetting, totalAmount, payload.ReceiptPartyType, receipt.Code,
		receipt.CostCenterID, receipt.ProjectID,receipt.PostingTime,receipt.PostingDate,
	payload.CompanyDefault.Currency)
	err = r.transactioRepo.SaveTransactionsBatch(ctx, tx, transactions)
	return err
}

func (r *atxBuyingEventRepo) transactionInventory(stockSetting model.StockSetting,
	lineItemAmount int64, voucherType, voucherCode string, costCenterID, projectID *int64,
	postingTime string,postingDate time.Time,currency string) []*model.TransactionLedger {
	credit := model.TransactionLedger{
		Ledger:         stockSetting.InventoryAccount,
		LedgerAgainst:  &stockSetting.StockReceivedButNotBilled,
		Debit:          lineItemAmount,
		VoucherCode:    voucherCode,
		VoucherType:    voucherType,
		VoucherSubtype: voucherType,
		ProjectID:      projectID,
		CostCenterID:   costCenterID,
		PostingTime: postingTime,
		PostingDate: postingDate,
		Currency: currency,
	}
	debit := model.TransactionLedger{
		Ledger:         stockSetting.StockReceivedButNotBilled,
		LedgerAgainst:  &stockSetting.InventoryAccount,
		Credit:         lineItemAmount,
		VoucherCode:    voucherCode,
		VoucherType:    voucherType,
		VoucherSubtype: voucherType,
		ProjectID:      projectID,
		CostCenterID:   costCenterID,
		PostingTime: postingTime,
		PostingDate: postingDate,
		Currency: currency,
	}
	return []*model.TransactionLedger{&credit, &debit}
}
