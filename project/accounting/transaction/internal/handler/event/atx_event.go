package transaction_event

import (
	"context"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/logger"
	transaction_repo "erp/project/accounting/transaction/internal/repository"
)

type atxEventHandler struct {
	emitLog           logger.EmitLog
	atxBuyinhRepo     transaction_repo.AtxBuyingEventRepo
	atxSellingRepo    transaction_repo.AtxSellingEventRepo
	atxStockEntryRepo transaction_repo.AtxStockEntryEventRepo
	accounting repository.AccountingService
}

func NewAtxEBuyingventHandler(
	bus bus.Bus,
	logger logger.Logger,
	atxBuyingRepo transaction_repo.AtxBuyingEventRepo,
	atxSellingRepo transaction_repo.AtxSellingEventRepo,
	atxStockEntryRepo transaction_repo.AtxStockEntryEventRepo,
	accounting repository.AccountingService,
) {

	handler := atxEventHandler{
		emitLog:           logger.EmitLog("atx-event"),
		atxBuyinhRepo:     atxBuyingRepo,
		atxSellingRepo:    atxSellingRepo,
		atxStockEntryRepo: atxStockEntryRepo,
		accounting: accounting,
	}
	bus.RegisterHandler(domain.ReceiptSubmittedEvent, handler.OnReceiptSubmitted())
	bus.RegisterHandler(domain.StockEntrySubmittedEvent, handler.OnStockEntrySubmitted())
	bus.RegisterHandler(domain.InvoiceSubmittedEvent, handler.OnInvoiceSubmitted())
	bus.RegisterHandler(domain.InvoiceCancelledEvent, handler.OnInvoiceCancelled())
	bus.RegisterHandler(domain.ReceiptCancelledEvent, handler.OnReceiptCancelled())
}

func (h *atxEventHandler) OnReceiptCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			paylaod,ok := e.Data.(event.StatusReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.accounting.DelTxnsByVoucherCode(ctx,paylaod.Tx,paylaod.Receipt.Code)
			return
		},
		AbortOnError: true,
		Matcher: domain.ReceiptCancelledEvent,
	}
}

func (h *atxEventHandler) OnInvoiceCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			paylaod,ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.accounting.DelTxnsByVoucherCode(ctx,paylaod.Tx,paylaod.Invoice.Code)
			return
		},
		AbortOnError: true,
		Matcher: domain.InvoiceCancelledEvent,
	}
}

func (h *atxEventHandler) OnStockEntrySubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnStockEntrySubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusStockEntryEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.atxStockEntryRepo.OnStockEntrySubmitted(ctx, payload)
			return err
		},
		AbortOnError: true,
		Matcher:      domain.StockEntrySubmittedEvent,
	}
}

func (h *atxEventHandler) OnInvoiceSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnInvoiceSubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			switch payload.InvoicePartyType {
			case proto.PartyType_purchaseInvoice.String():
				err = h.atxBuyinhRepo.OnInvoiceSubmitted(ctx, payload)
			case proto.PartyType_saleInvoice.String():
				err = h.atxSellingRepo.OnInvoiceSubmitted(ctx, payload)
			}
			return err
		},
		AbortOnError: true,
		Matcher:      domain.InvoiceSubmittedEvent,
	}
}

func (h *atxEventHandler) OnReceiptSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnReceiptSubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			switch payload.ReceiptPartyType {
			case proto.PartyType_purchaseReceipt.String():
				err = h.atxBuyinhRepo.OnReceiptSubmitted(ctx, payload)
			case proto.PartyType_deliveryNote.String():
				err = h.atxSellingRepo.OnDeleveryNoteSubmitted(ctx, payload)
			}
			return err
		},
		AbortOnError: true,
		Matcher:      domain.ReceiptSubmittedEvent,
	}
}
