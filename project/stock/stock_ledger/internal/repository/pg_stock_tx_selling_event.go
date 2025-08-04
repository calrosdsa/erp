package stock_ledger_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"

	// "erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"fmt"

	"gorm.io/gorm"
)

type StockTxSellingRepository interface {
	OnDeliveryNoteSubmitted(ctx context.Context, payload event.StatusReceiptEventData) error
	OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) error
}

type stockTxSellingRepo struct {
	stockLedgerRepo StockLedgerTxRepository
	accounting     repository.AccountingService
	currency helpers.CurrencyHelper
}

func NewStockTxSellingRepo(
	stockLedgerRepo StockLedgerTxRepository,
	helpers *helpers.Helpers,
	accounting repository.AccountingService,
) StockTxSellingRepository {
	return &stockTxSellingRepo{
		stockLedgerRepo: stockLedgerRepo,
		accounting: accounting,
		currency: helpers.Currency,
	}
}

func (r *stockTxSellingRepo) OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) (err error) {
	if !payload.Invoice.UpdateStock {
		return nil
	}
	tx := payload.Tx
	invoice := payload.Invoice
	companyDefault := payload.CompanyDefaults
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefault,invoice.Currency, true, false)
	stockTxs := make([]*model.StockTransaction, len(payload.LineItemsData.LineItems))
	for i, line := range payload.LineItemsData.LineItems {
		if !line.MaintainStock {
			return
		}
		rate := r.currency.CurrencyExchange(int64(line.Rate), exchangeRate)
		stockLedgerBalance, err := r.getBalanceItem(ctx, tx, line.ItemID,
			line.SourceWarehouseID, rate, -line.Quantity)
		if err != nil {
			return err
		}
		fmt.Println("STOCK BALANCE INVOICE", stockLedgerBalance)
		stockTx := &model.StockTransaction{}
		stockTx.ItemID = line.ItemID
		stockTx.UomID = line.UnitOfMeasureID
		stockTx.OutQty = line.Quantity
		stockTx.BalanceQty = stockLedgerBalance.BalanceQuantity
		stockTx.WarehouseID = line.SourceWarehouseID
		stockTx.IncomingRate = rate
		stockTx.AverageRate = stockLedgerBalance.AvgRate
		stockTx.ValuationRate = rate
		stockTx.PostingDate = invoice.PostingDate
		stockTx.BalanceValue = stockLedgerBalance.BalanceRate
		stockTx.VoucherType = proto.PartyType_deliveryNote.String()
		stockTx.VoucherNo = invoice.Code
		stockTx.Currency = payload.Invoice.Currency
		stockTxs[i] = stockTx
	}
	err = tx.StockTransaction.WithContext(ctx).CreateInBatches(stockTxs, len(stockTxs))
	if err != nil {
		return
	}
	return
}

func (r *stockTxSellingRepo) OnDeliveryNoteSubmitted(ctx context.Context,
	payload event.StatusReceiptEventData) (err error) {
	tx := payload.Tx
	companyDefault := payload.CompanyDefault
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefault,payload.Receipt.Currency, true, false)
	
	fmt.Println("ON DELIVERY NOTE SUBMITTED...")
	stockTxs := make([]*model.StockTransaction, len(payload.LineItemsData.LineItems))

	if err != nil {
		return
	}
	for i, line := range payload.LineItemsData.LineItems {
		if !line.MaintainStock {
			continue
		}
		rate := r.currency.CurrencyExchange(int64(line.Rate), exchangeRate)
		stockLedgerBalance, err := r.getBalanceItem(ctx, tx, line.ItemID,
			line.SourceWarehouseID, rate, -line.Quantity)
		if err != nil {
			return err
		}

		fmt.Println("STOCK BALANCE", stockLedgerBalance)
		stockTx := &model.StockTransaction{}
		stockTx.ItemID = line.ItemID
		stockTx.UomID = line.UnitOfMeasureID
		stockTx.OutQty = line.Quantity
		stockTx.BalanceQty = stockLedgerBalance.BalanceQuantity
		stockTx.WarehouseID = line.SourceWarehouseID
		stockTx.IncomingRate = rate
		stockTx.AverageRate = stockLedgerBalance.AvgRate
		stockTx.ValuationRate = rate
		stockTx.PostingDate = payload.Receipt.PostingDate
		stockTx.BalanceValue = stockLedgerBalance.BalanceRate
		stockTx.VoucherType = proto.PartyType_deliveryNote.String()
		stockTx.VoucherNo = payload.Receipt.Code
		stockTx.Currency = payload.Receipt.Currency
		stockTxs[i] = stockTx
		// err = r.reduceStockQty(ctx, tx, itemLine.AcceptedQuantity,
		// 	itemPrice.ItemID, itemLine.AcceptedWarehouse, payload.StockDefault)
		// if err!= nil {
		// 	return err
		// }
	}
	err = tx.StockTransaction.WithContext(ctx).CreateInBatches(stockTxs, len(stockTxs))
	if err != nil {
		return
	}
	return err
}

func (r *stockTxSellingRepo) getBalanceItem(ctx context.Context, tx *query.QueryTx, args ...interface{}) (res StockBalance, err error) {
	itemID := args[0].(int64)
	warehouseID := args[1].(int64)
	incomingRate := args[2].(int64)
	in_qty := args[3].(int32)
	stockTx := tx.StockTransaction
	builder := stockTx.WithContext(ctx).Where(
		stockTx.ItemID.Eq(itemID),
		stockTx.WarehouseID.Eq(warehouseID),
	)
	count, err := builder.Count()
	if err != nil {
		fmt.Println("GET BALANCE ITEM", err)
	}
	//get last in stock transaction
	inStockTx, err := stockTx.WithContext(ctx).Where(
		stockTx.ItemID.Eq(itemID),
		stockTx.WarehouseID.Eq(warehouseID),
		stockTx.OutQty.Eq(0),
	).Select(
		stockTx.IncomingRate,
	).Where().Last()
	if err == gorm.ErrRecordNotFound {
		err = nil
		inStockTx = &model.StockTransaction{}
		inStockTx.IncomingRate = incomingRate
	}
	if err != nil {
		return
	}
	//get last register
	stockTransaction, err := builder.
		Select(stockTx.BalanceQty, stockTx.BalanceValue, stockTx.AverageRate).
		Last()
	if err == gorm.ErrRecordNotFound {
		err = nil
		stockTransaction = &model.StockTransaction{}
	}
	if err != nil {
		return
	}
	fmt.Println("COUNT", count)
	fmt.Println("STOCK TX", stockTransaction)
	res.AvgRate = ((stockTransaction.AverageRate * count) + incomingRate) / (count+1)
	res.BalanceQuantity = stockTransaction.BalanceQty + in_qty
	res.BalanceRate = stockTransaction.BalanceValue + (int64(in_qty) * inStockTx.IncomingRate)

	return
}
