package tac_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	tac_repo "erp/project/accounting/tax_and_charges/repository"
)

type tacEventHandler struct {
	emiLog       logger.EmitLog
	tacEventRepo tac_repo.TacRepositoryEvent
}

func NewTacEventHandler(
	bus bus.Bus,
	tacEventRepo tac_repo.TacRepositoryEvent,
	logger logger.Logger,
) {
	handler := tacEventHandler{
		emiLog:       logger.EmitLog("tac-event"),
		tacEventRepo: tacEventRepo,
	}
	bus.RegisterHandler(domain.QuotationCreatedEvent, handler.OnQuotationCreated())
	bus.RegisterHandler(domain.ReceiptCreatedEvent, handler.OnReceiptCreated())
	bus.RegisterHandler(domain.OrderCreatedEvent, handler.OnOrderCreated())
	bus.RegisterHandler(domain.InvoiceCreatedEvent, handler.OnInvoiceCreated())
	bus.RegisterHandler(domain.CashOutflowCreatedEvent,handler.OnCashOutflowCreated())

	bus.RegisterHandler(domain.QuotationEditEvent, handler.OnQuotationEdited())
	bus.RegisterHandler(domain.ReceiptEditEvent, handler.OnReceiptEdited())
	bus.RegisterHandler(domain.InvoiceEditEvent, handler.OnInvoiceEdited())
	bus.RegisterHandler(domain.OrderEditEvent, handler.OnOrderEdited())
	bus.RegisterHandler(domain.CashOutflowEditedEvent,handler.OnCashOutflowEdited())

	bus.RegisterHandler(domain.ChargesTemplateCreatedEvent, handler.OnChargesTemplateCreated())
	bus.RegisterHandler(domain.PaymentCreatedEvent, handler.OnPaymentCreated())
	bus.RegisterHandler(domain.PaymentEditedEvent, handler.OnPaymentEdited())

}
func (h *tacEventHandler) OnPaymentEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnPaymentEdited"))
				}
			}()
			payload, ok := e.Data.(event.PaymentEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.EditTaxAndChargeLines(payload.Tx, ctx, payload.Body.CreateTaxAndCharges, payload.Body.PaymentData.ID)
			return
		},
		AbortOnError: true,
		Matcher:      domain.PaymentEditedEvent,
	}
}

func (h *tacEventHandler) OnReceiptEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnReceiptEdited"))
				}
			}()
			payload, ok := e.Data.(event.ReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.EditTaxAndChargeLines(payload.Tx, ctx, payload.Body.CreateTaxAndCharges, payload.Body.Receipt.ID)
			return
		},
		AbortOnError: true,
		Matcher:      domain.ReceiptEditEvent,
	}
}

func (h *tacEventHandler) OnQuotationEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnQuotationEdited"))
				}
			}()
			payload, ok := e.Data.(event.QuotationEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.EditTaxAndChargeLines(payload.Tx, ctx, payload.Body.CreateTaxAndCharges, payload.Body.Quotation.ID)
			return
		},
		AbortOnError: true,
		Matcher:      domain.QuotationEditEvent,
	}
}

func (h *tacEventHandler) OnOrderEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnOrderEdited"))
				}
			}()
			payload, ok := e.Data.(event.OrderEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.EditTaxAndChargeLines(payload.Tx, ctx, payload.Body.CreateTaxAndCharges, payload.Body.Order.ID)
			return
		},
		AbortOnError: true,
		Matcher:      domain.OrderEditEvent,
	}
}

func (h *tacEventHandler) OnInvoiceEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnInvoiceEdited"))
				}
			}()
			payload, ok := e.Data.(event.InvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.EditTaxAndChargeLines(payload.Tx, ctx, payload.Body.CreateTaxAndCharges, payload.Body.Invoice.ID)
			return
		},
		AbortOnError: true,
		Matcher:      domain.InvoiceEditEvent,
	}
}


func (h *tacEventHandler) OnCashOutflowEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnCashOutflowEdited"))
				}
			}()
			payload, ok := e.Data.(event.CashOutflowEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.EditTaxAndChargeLines(
				payload.Tx, ctx, payload.Data.CreateTaxAndCharges,
				payload.Data.ID,)
			return
		},
		AbortOnError: true,
		Matcher:      domain.CashOutflowEditedEvent,
	}
}

func (h *tacEventHandler) OnCashOutflowCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnCashOutflowCreated"))
				}
			}()
			payload, ok := e.Data.(event.CashOutflowEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.CreateTaxAndChargeLines(
				payload.Tx, ctx, payload.Data.CreateTaxAndCharges,
				payload.CashOutflow.ID)
			return
		},
		AbortOnError: true,
		Matcher:      domain.CashOutflowCreatedEvent,
	}
}

func (h *tacEventHandler) OnPaymentCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnPaymentCreated"))
				}
			}()
			payload, ok := e.Data.(event.PaymentEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.OnPaymentCreated(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.PaymentCreatedEvent,
	}
}

func (h *tacEventHandler) OnChargesTemplateCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnChargesTemplateCreated"))
				}
			}()
			payload, ok := e.Data.(event.ChargesTemplateEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.OnChargesTemplateCreated(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.ChargesTemplateCreatedEvent,
	}
}

func (h *tacEventHandler) OnReceiptCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnReceiptCreated"))
				}
			}()
			payload, ok := e.Data.(event.ReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.OnReceiptCreated(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.ReceiptCreatedEvent,
	}
}

func (h *tacEventHandler) OnOrderCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnOrderCreated"))
				}
			}()
			payload, ok := e.Data.(event.OrderEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.OnOrderCreated(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.OrderCreatedEvent,
	}
}

func (h *tacEventHandler) OnInvoiceCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnInvoiceCreated"))
				}
			}()
			payload, ok := e.Data.(event.InvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.OnInvoiceCreated(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.InvoiceCreatedEvent,
	}
}

func (h *tacEventHandler) OnQuotationCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emiLog.Err(err, logger.OptionsLog.WithMethod("OnQuotationCreated"))
				}
			}()
			payload, ok := e.Data.(event.QuotationEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.tacEventRepo.OnQuotationCreated(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.QuotationCreatedEvent,
	}
}
