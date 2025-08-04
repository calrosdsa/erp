package stock_ledger_event

import (
	"context"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	stock_ledger_repo "erp/project/stock/stock_ledger/internal/repository"
	"fmt"
)

type stockTxBuyingHandler struct {
	emitLog            logger.EmitLog
	stockTxBuyingRepo  stock_ledger_repo.StockTxBuyingRepository
	stockTxSellingRepo stock_ledger_repo.StockTxSellingRepository
	stockTxEntry       stock_ledger_repo.StockTxStockEntry
	stockTxRepo        stock_ledger_repo.StockLedgerTxRepository
}

func NewStockTxBuyingHandler(
	bus bus.Bus,
	logger logger.Logger,
	stockTxBuyingRepo stock_ledger_repo.StockTxBuyingRepository,
	stockTxSellingRepo stock_ledger_repo.StockTxSellingRepository,
	stockTxEntry stock_ledger_repo.StockTxStockEntry,
	stockTxRepo stock_ledger_repo.StockLedgerTxRepository,
) {
	h := stockTxBuyingHandler{
		emitLog:            logger.EmitLog("stock-ledger"),
		stockTxBuyingRepo:  stockTxBuyingRepo,
		stockTxSellingRepo: stockTxSellingRepo,
		stockTxEntry:       stockTxEntry,
		stockTxRepo:        stockTxRepo,
	}
	fmt.Println("STOCK LEDGER EVENT REGISTERS")
	bus.RegisterHandler(domain.ReceiptSubmittedEvent, h.OnReceiptSubmitted())
	bus.RegisterHandler(domain.StockEntrySubmittedEvent, h.OnStockEntrySubmitted())
	bus.RegisterHandler(domain.InvoiceSubmittedEvent, h.OnInvoiceSubmitted())
	bus.RegisterHandler(domain.InvoiceCancelledEvent, h.OnInvoiceCancelled())
	bus.RegisterHandler(domain.ReceiptCancelledEvent, h.OnReceiptCancelled())
}
func (h *stockTxBuyingHandler) OnReceiptCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnReceiptCancelled"))
				}
			}()
			payload, ok := e.Data.(event.StatusReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.stockTxRepo.OnCancelled(payload.Tx,ctx,payload.Receipt.Code)
			return
		},
		AbortOnError: true,
		Matcher:      domain.ReceiptCancelledEvent,
	}
}


func (h *stockTxBuyingHandler) OnInvoiceCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnInvoiceCancelled"))
				}
			}()
			payload, ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.stockTxRepo.OnCancelled(payload.Tx,ctx,payload.Invoice.Code)
			return
		},
		AbortOnError: true,
		Matcher:      domain.InvoiceCancelledEvent,
	}
}


func (h *stockTxBuyingHandler) OnInvoiceSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnInvoiceSubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			switch payload.InvoicePartyType {
			case proto.PartyType_purchaseInvoice.String():
				err = h.stockTxBuyingRepo.OnInvoiceSubmitted(ctx, payload)
			case proto.PartyType_saleInvoice.String():
				err = h.stockTxSellingRepo.OnInvoiceSubmitted(ctx, payload)
			}
			return
		},
		AbortOnError: true,
		Matcher:      domain.InvoiceSubmittedEvent,
	}
}

func (h *stockTxBuyingHandler) OnStockEntrySubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnStockEntrySubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusStockEntryEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.stockTxEntry.OnStockEntrySubmitted(ctx, payload)
			return err
		},
		AbortOnError: true,
		Matcher:      domain.StockEntrySubmittedEvent,
	}
}

func (h *stockTxBuyingHandler) OnReceiptSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnReceiptSubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			switch payload.ReceiptPartyType {
			case proto.PartyType_purchaseReceipt.String():
				err = h.stockTxBuyingRepo.OnReceiptSubmitted(ctx, payload)
			case proto.PartyType_deliveryNote.String():
				err = h.stockTxSellingRepo.OnDeliveryNoteSubmitted(ctx, payload)
			}
			return err
		},
		AbortOnError: true,
		Matcher:      domain.ReceiptSubmittedEvent,
	}
}
