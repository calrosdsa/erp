package invoice_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	invoice_repo "erp/project/invoice/internal/repository"
)

type InvoiceEventHandler struct {
	emitLog          logger.EmitLog
	invoiceEventRepo invoice_repo.InvoiceEventRepository
}

func NewInvoiceEventHandler(
	bus bus.Bus,
	logger logger.Logger,
	invoiceEventRepo invoice_repo.InvoiceEventRepository,
) {
	handler := InvoiceEventHandler{
		emitLog:          logger.EmitLog("invoice-event-handler"),
		invoiceEventRepo: invoiceEventRepo,
	}
	bus.RegisterHandler(domain.ReceiptSubmittedEvent, handler.OnReceiptSubmitted())
	bus.RegisterHandler(domain.PaymentSubmittedEvent, handler.OnPaymentSubmitted())
	bus.RegisterHandler(domain.PaymentCancelledEvent, handler.OnPaymentCancelled())
}
func (h *InvoiceEventHandler) OnPaymentCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.StatusPaymentEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.invoiceEventRepo.OnPaymentCancelled(ctx, payload)
			if err != nil {
				h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnPaymentCancelled"))
				return err
			}
			return nil
		},
		AbortOnError: true,
		Matcher:      domain.PaymentCancelledEvent,
	}
}
func (h *InvoiceEventHandler) OnPaymentSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.StatusPaymentEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.invoiceEventRepo.OnPaymentSubmitted(ctx, payload)
			if err != nil {
				h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnPaymentSubmitted"))
				return err
			}
			return nil
		},
		Matcher:      domain.PaymentSubmittedEvent,
		AbortOnError: true,
	}
}

func (h *InvoiceEventHandler) OnReceiptSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.StatusReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.invoiceEventRepo.OnReceiptSubmitted(ctx, payload)
			if err != nil {
				h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnReceiptSubmitted"))
				return err
			}
			return nil
		},
		Matcher:      domain.ReceiptSubmittedEvent,
		AbortOnError: true,
	}
}
