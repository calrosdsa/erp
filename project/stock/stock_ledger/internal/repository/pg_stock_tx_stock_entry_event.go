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

type StockTxStockEntry interface {
	OnStockEntrySubmitted(ctx context.Context, payload event.StatusStockEntryEventData) (err error)
}

type stockTxStockEntry struct {
	stockLedgerRepo StockLedgerTxRepository
	currency helpers.CurrencyHelper
	accounting     repository.AccountingService
}

func NewStockTxStockEntryRepo(
	stockLedgerRepo StockLedgerTxRepository,
	helpers *helpers.Helpers,
	accounting     repository.AccountingService,
) StockTxStockEntry {
	return &stockTxStockEntry{
		stockLedgerRepo: stockLedgerRepo,
		currency: helpers.Currency,
		accounting: accounting,
	}
}

func (r *stockTxStockEntry) OnStockEntrySubmitted(ctx context.Context, payload event.StatusStockEntryEventData) (err error) {
	tx := payload.Tx
	stockTxs := make([]*model.StockTransaction, len(payload.LineItemsData.LineItems))
	for i, line := range payload.LineItemsData.LineItems {
		var warehouseID int64
		var stockBalance StockBalance
		if !line.MaintainStock {
			continue
		}
		switch payload.StockEntry.EntryType {
		case proto.StockEntryType_MATERIAL_RECEIPT.String():
			if line.TargetWarehouseID != 0 {
				warehouseID = line.TargetWarehouseID
				stockBalance, err = r.stockLedgerRepo.GetBalanceItem(ctx, tx, line.ItemID,
					warehouseID, line.Rate, line.Quantity)
			}
		}
		if err != nil {
			return err
		}
		fmt.Println("STOCK BALANCE", stockBalance)
		stockTx := &model.StockTransaction{}
		stockTx.ItemID = line.ItemID
		stockTx.UomID = line.UnitOfMeasureID
		stockTx.InQty = line.Quantity
		stockTx.BalanceQty = stockBalance.BalanceQuantity
		stockTx.WarehouseID = warehouseID
		stockTx.IncomingRate = line.Rate
		stockTx.AverageRate = stockBalance.AvgRate
		stockTx.ValuationRate = line.Rate
		stockTx.AvailableQty = line.Quantity
		stockTx.BalanceValue = stockBalance.BalanceRate
		stockTx.PostingDate = payload.StockEntry.PostingDate
		stockTx.VoucherType = proto.PartyType_stockEntry.String()
		stockTx.VoucherNo = payload.StockEntry.Code
		stockTx.Currency = payload.CompanyDefaults.Currency
		stockTxs[i] = stockTx

	}
	err = tx.StockTransaction.WithContext(ctx).CreateInBatches(stockTxs, len(stockTxs))
	return err
}
