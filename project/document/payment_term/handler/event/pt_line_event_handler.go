package payment_terms_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	payment_terms_repo "erp/project/document/payment_term/repository"
)

type ptLineEventHandler struct {
	emitLog    logger.EmitLog
	ptLineRepo payment_terms_repo.PaymentTermsLineRepo
}

func NewPtLineEventHandler(
	logger logger.Logger,
	bus bus.Bus,
	ptLineRepo payment_terms_repo.PaymentTermsLineRepo,
) {
	h := ptLineEventHandler{
		emitLog:    logger.EmitLog("pt-line-event-handler"),
		ptLineRepo: ptLineRepo,
	}

	bus.RegisterHandler(domain.PaymentTermsTemplateCreatedEvent,h.OnPaymentTemplateCreated())
	bus.RegisterHandler(domain.PaymentTermsTemplateEditedEvent,h.OnPaymentTemplateEdited())

}

func (h *ptLineEventHandler) OnPaymentTemplateCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnPaymentTemplateCreated"))
				}
			}()
			payload, ok := e.Data.(event.PaymentTermsTemplateEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.ptLineRepo.CreatePaymentTermLines(payload.Tx, ctx, payload.Body.Lines, payload.PaymentTermsTemplateID)
			return err
		},
		AbortOnError: true,
		Matcher: domain.PaymentTermsTemplateCreatedEvent,
	}
}

func (h *ptLineEventHandler) OnPaymentTemplateEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnPaymentTemplateEdited"))
				}
			}()
			payload, ok := e.Data.(event.PaymentTermsTemplateEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.ptLineRepo.UpdatePaymentTermsLines(payload.Tx, ctx, payload.Body.Lines, payload.PaymentTermsTemplateID)
			return err
		},
		AbortOnError: true,
		Matcher: domain.PaymentTermsTemplateEditedEvent,
	}
}
