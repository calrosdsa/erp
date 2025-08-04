package receipt_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	receipt_repo "erp/project/receipt/internal/repository"
)

type ReceiptEventHandler struct {
	emitLog          logger.EmitLog
	receiptEventRepo receipt_repo.ReceiptEventRepository
}

func NewReceiptEventHandler(
	bus bus.Bus,
	logger logger.Logger,
	receiptEventRepo receipt_repo.ReceiptEventRepository,
) {
	handler := ReceiptEventHandler{
		emitLog:          logger.EmitLog("receipt-event-handler"),
		receiptEventRepo: receiptEventRepo,
	}
	bus.RegisterHandler(domain.InvoiceSubmittedEvent, handler.OnInvoiceSubmitted())
}

func (h *ReceiptEventHandler) OnInvoiceSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.receiptEventRepo.OnInvoiceSubmitted(ctx, payload)
			if err != nil {
				h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnInvoiceSubmitted"))
				return err
			}
			return nil
		},
		Matcher:      domain.InvoiceSubmittedEvent,
		AbortOnError: true,
	}
}
