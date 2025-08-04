package serial_no_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	serial_no_repo "erp/project/stock/serial_no/internal/repository"
	"fmt"
)

type serialNoEventHandler struct {
	emitLog           logger.EmitLog
	serialNoEventRepo serial_no_repo.SerialNoEventRepository
}

func NewSerialNoEventHandler(
	logger logger.Logger,
	bus bus.Bus,
	serialNoEventRepo serial_no_repo.SerialNoEventRepository,
) {
	handler := serialNoEventHandler{
		emitLog:           logger.EmitLog("serial-no-event"),
		serialNoEventRepo: serialNoEventRepo,
	}
	fmt.Println("REGISTER SERIAL EVENT REPO")
	bus.RegisterHandler(domain.InvoiceSubmittedEvent, handler.OnInvoiceSubmitted())
	bus.RegisterHandler(domain.ReceiptSubmittedEvent, handler.OnReceiptSubmitted())
	bus.RegisterHandler(domain.StockEntrySubmittedEvent, handler.OnStockEntrySubmitted())
	bus.RegisterHandler(domain.InvoiceCancelledEvent, handler.OnInvoiceCancelled())
	bus.RegisterHandler(domain.ReceiptCancelledEvent, handler.OnReceiptCancelled())
}

func (h *serialNoEventHandler) OnStockEntrySubmitted() bus.Handler {
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
			err = h.serialNoEventRepo.OnStockEntrySubmitted(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.StockEntrySubmittedEvent,
	}
}

func (h *serialNoEventHandler) OnReceiptSubmitted() bus.Handler {
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
			err = h.serialNoEventRepo.OnReceiptSubmitted(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.ReceiptSubmittedEvent,
	}
}

func (h *serialNoEventHandler) OnReceiptCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnReceiptCancelled"))
				}
			}()
			payload, ok := e.Data.(event.StatusReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.serialNoEventRepo.OnReceiptCancelled(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.ReceiptCancelledEvent,
	}
}

func (h *serialNoEventHandler) OnInvoiceCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnInvoiceCancelled"))
				}
			}()
			payload, ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.serialNoEventRepo.OnInvoiceCancelled(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.InvoiceCancelledEvent,
	}
}

func (h *serialNoEventHandler) OnInvoiceSubmitted() bus.Handler {
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
			err = h.serialNoEventRepo.OnInvoiceSubmitted(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.InvoiceSubmittedEvent,
	}
}
