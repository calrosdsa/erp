package payment_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	payment_repo "erp/project/accounting/payment/internal/repository"
)

type paymentEventHandler struct {
	bus bus.Bus
	emiLog logger.EmitLog
	paymentEventRepo payment_repo.PaymentEventRepository
}

func NewPaymentEventHandler(
	bus bus.Bus,
	logger logger.Logger,
	paymentEventRepo payment_repo.PaymentEventRepository,
) {
	h := paymentEventHandler{
		emiLog: logger.EmitLog("payment-event-handler"),		
		paymentEventRepo: paymentEventRepo,
	}
	bus.RegisterHandler(domain.InvoiceCancelledEvent,h.OnInvoiceCancelled())
}
func (h *paymentEventHandler)OnInvoiceCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emiLog.Err(err,logger.OptionsLog.WithMethod("OnInvoiceCancelled"))
				}
			}()
			payload,ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				 return domain.FAIL_TYPE_ASSERTION
			}
			err = h.paymentEventRepo.OnInvoiceCancelled(ctx,payload)
			if err != nil {
				return
			}
			return nil
		},
		AbortOnError: true,
		Matcher: domain.InvoiceCancelledEvent,
	}
}