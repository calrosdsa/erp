package stock_ledger_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"fmt"
)

type StockTxBuyingRepository interface {
	OnReceiptSubmitted(ctx context.Context, payload event.StatusReceiptEventData) error
	OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) error
}

type stockTxBuyingRepo struct {
	stockLedgerRepo StockLedgerTxRepository
	accounting     repository.AccountingService
	currency       helpers.CurrencyHelper
}

func NewStockTxBuyingRepo(
	stockLedgerRepo StockLedgerTxRepository,
	accounting     repository.AccountingService,
	helpers       *helpers.Helpers,
) StockTxBuyingRepository {
	return &stockTxBuyingRepo{
		stockLedgerRepo: stockLedgerRepo,
		accounting:accounting,
		currency: helpers.Currency,
	}
}

func (r *stockTxBuyingRepo) OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) error {
	if !payload.Invoice.UpdateStock {
		return nil
	}
	tx := payload.Tx
	invoice := payload.Invoice
	companyDefault := payload.CompanyDefaults
	
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefault,invoice.Currency, false, true)
	if err != nil  {
		return err
	}
	stockTxs := make([]*model.StockTransaction, len(payload.LineItemsData.LineItems))
	for i, line := range payload.LineItemsData.LineItems {
		fmt.Println("STOCK TX  OnInvoiceSubmitted LINE",line)
		if !line.MaintainStock {
			continue
		}
		rate := r.currency.CurrencyExchange(int64(line.Rate), exchangeRate)
		stockLedgerBalance, err := r.stockLedgerRepo.GetBalanceItem(ctx, tx, line.ItemID,
			line.AcceptedWarehouse, rate, line.AcceptedQuantity)
		if err != nil {
			return err
		}
		stockTx := &model.StockTransaction{}
		stockTx.ItemID = line.ItemID
		stockTx.UomID = line.UnitOfMeasureID
		stockTx.InQty = line.AcceptedQuantity
		stockTx.BalanceQty = stockLedgerBalance.BalanceQuantity
		stockTx.WarehouseID = line.AcceptedWarehouse
		stockTx.IncomingRate = rate
		stockTx.AverageRate = stockLedgerBalance.AvgRate
		stockTx.ValuationRate = rate
		stockTx.PostingDate = payload.Invoice.PostingDate
		stockTx.AvailableQty = line.AcceptedQuantity
		stockTx.BalanceValue = stockLedgerBalance.BalanceRate
		stockTx.VoucherType = proto.PartyType_purchaseInvoice.String()
		stockTx.VoucherNo = payload.Invoice.Code
		stockTx.Currency = companyDefault.Currency
		stockTxs[i] = stockTx
	}
	err = tx.StockTransaction.WithContext(ctx).CreateInBatches(stockTxs, len(stockTxs))
	return err
}

func (r *stockTxBuyingRepo) OnReceiptSubmitted(ctx context.Context,
	payload event.StatusReceiptEventData) (err error) {
	tx := payload.Tx
	receipt := payload.Receipt
	companyDefault := payload.CompanyDefault
	stockTxs := make([]*model.StockTransaction, len(payload.LineItemsData.LineItems))
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefault, receipt.Currency, false, true)
	if err != nil {
		return
	}
	for i, line := range payload.LineItemsData.LineItems {
		if !line.MaintainStock {
			return
		}
		rate := r.currency.CurrencyExchange(int64(line.Rate), exchangeRate)
		stockLedgerBalance, err := r.stockLedgerRepo.GetBalanceItem(ctx, tx, line.ItemID,
			line.AcceptedWarehouse, rate, line.Quantity)
		if err != nil {
			return err
		}
		fmt.Println("STOCK BALANCE", stockLedgerBalance)

		stockTx := &model.StockTransaction{}
		stockTx.ItemID = line.ItemID
		stockTx.UomID = line.UnitOfMeasureID
		stockTx.InQty = line.Quantity
		stockTx.BalanceQty = stockLedgerBalance.BalanceQuantity
		stockTx.WarehouseID = line.AcceptedWarehouse
		stockTx.IncomingRate = rate
		stockTx.AverageRate = stockLedgerBalance.AvgRate
		stockTx.ValuationRate = rate
		stockTx.PostingDate = payload.Receipt.PostingDate
		stockTx.AvailableQty = line.Quantity
		stockTx.BalanceValue = stockLedgerBalance.BalanceRate
		stockTx.VoucherType = proto.PartyType_purchaseReceipt.String()
		stockTx.VoucherNo = payload.Receipt.Code
		stockTx.Currency = companyDefault.Currency
		stockTxs[i] = stockTx
	}
	err = tx.StockTransaction.WithContext(ctx).CreateInBatches(stockTxs, len(stockTxs))
	return err
}

type ItemLineReceipt struct {
	model.ItemLine
	AcceptedQuantity  int32
	RejectedQuantity  int32
	AcceptedWarehouse int64
}
