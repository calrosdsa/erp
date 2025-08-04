package transaction_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
)

type AtxStockEntryEventRepo interface {
	OnStockEntrySubmitted(ctx context.Context, payload event.StatusStockEntryEventData) (err error)
}
type atxStockEntryEventRepo struct {
	setting        repository.SettingService
	transactioRepo TransactionRepository
	// accounting     repository.AccountingService
	// currency       helpers.CurrencyHelper
}

func NewAtxStockEntryEventRepo(
	setting repository.SettingService,
	transactioRepo TransactionRepository,
	accounting     repository.AccountingService,
	helpers *helpers.Helpers,
) AtxStockEntryEventRepo {
	return &atxStockEntryEventRepo{
		setting:        setting,
		transactioRepo: transactioRepo,
		// accounting: accounting,
		// currency: helpers.Currency,
	}
}

func (r *atxStockEntryEventRepo) OnStockEntrySubmitted(ctx context.Context,
	payload event.StatusStockEntryEventData) (err error) {
	tx := payload.Tx
	stockEntry := payload.StockEntry
	amountItemLines := payload.LineItemsData.TotalAmount
	// Return if no amount is involved in the transactions.
	if amountItemLines == 0 {
		return nil
	}
	stockSetting, err := r.setting.GetStockSettings(ctx, payload.CompanyID)
	if err != nil {
		return
	}
	
	accountSetting, err := r.setting.GetAccountSettings(ctx, payload.CompanyID)
	if err != nil {
		return
	}
	//Adding receipt account entry
	debitTx := model.TransactionLedger{
		VoucherCode:    stockEntry.Code,
		VoucherType:    proto.PartyType_stockEntry.String(),
		VoucherSubtype: stockEntry.EntryType,
		PostingDate: stockEntry.PostingDate,
		PostingTime: stockEntry.PostingTime,
		Currency: payload.CompanyDefaults.Currency,
	}
	creditTx := model.TransactionLedger{
		VoucherCode:    stockEntry.Code,
		VoucherType:    proto.PartyType_stockEntry.String(),
		VoucherSubtype: stockEntry.EntryType,
		PostingDate: stockEntry.PostingDate,
		PostingTime: stockEntry.PostingTime,
		Currency: payload.CompanyDefaults.Currency,
	}
	switch payload.StockEntry.EntryType {
	case proto.StockEntryType_MATERIAL_ISSUE.String():
		debitTx.Ledger = accountSetting.CostOfGoodSoldAccount
		debitTx.LedgerAgainst = &stockSetting.InventoryAccount
		debitTx.Debit = amountItemLines

		creditTx.Ledger = stockSetting.InventoryAccount
		creditTx.LedgerAgainst = &accountSetting.CostOfGoodSoldAccount
		creditTx.Credit = amountItemLines
	case proto.StockEntryType_MATERIAL_RECEIPT.String():
		debitTx.Ledger = stockSetting.InventoryAccount
		debitTx.LedgerAgainst = &stockSetting.StockAdjustment
		debitTx.Debit = amountItemLines

		creditTx.Ledger = stockSetting.StockAdjustment
		creditTx.LedgerAgainst = &stockSetting.InventoryAccount
		creditTx.Credit = amountItemLines
	}
	err = r.transactioRepo.SaveTransactionsBatch(ctx, tx, []*model.TransactionLedger{&debitTx, &creditTx})
	return
}

