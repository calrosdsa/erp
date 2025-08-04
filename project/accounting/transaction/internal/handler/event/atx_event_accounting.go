package transaction_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	transaction_repo "erp/project/accounting/transaction/internal/repository"
)

type atxAccountingEventHandler struct {
	emitLog           logger.EmitLog
	atxAccountingRepo transaction_repo.AtxAccountignEventRepo
}

func NewAtxAccountingEventHandler(
	bus bus.Bus,
	logger logger.Logger,
	atxAccountRepo transaction_repo.AtxAccountignEventRepo,
) {
	h := atxAccountingEventHandler{
		emitLog:           logger.EmitLog("atx-accounting-event-handler"),
		atxAccountingRepo: atxAccountRepo,
	}
	bus.RegisterHandler(domain.PaymentSubmittedEvent, h.OnPaymentSubmitted())
	bus.RegisterHandler(domain.JournalEntrySubmittedEvent, h.OnJournalEntrySubmitted())
	bus.RegisterHandler(domain.JournalEntryCancelledEvent, h.OnJournalEntryCancelled())
	bus.RegisterHandler(domain.PaymentCancelledEvent, h.OnPaymentCancelled())

	bus.RegisterHandler(domain.CashOutflowCancelledEvent, h.OnCashOutflowCancelled())
	bus.RegisterHandler(domain.CashOutflowSubmittedEvent, h.OnCashOutflowSubmitted())
}

func (h *atxAccountingEventHandler) OnCashOutflowCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnCashOutflowCancelled"))
				}
			}()
			payload, ok := e.Data.(event.StatusCashOutflowEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.atxAccountingRepo.OnCashOutflowCancelled(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.CashOutflowCancelledEvent,
	}
}

func (h *atxAccountingEventHandler) OnCashOutflowSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnCashOutflowSubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusCashOutflowEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.atxAccountingRepo.OnCashOutflowSubmitted(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		},
		AbortOnError: true,
		Matcher:      domain.CashOutflowSubmittedEvent,
	}
}

func (h *atxAccountingEventHandler) OnJournalEntryCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnJournalEntryCancelled"))
				}
			}()
			payload, ok := e.Data.(event.StatusJournalEntryEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.atxAccountingRepo.OnJournalEntryCancelled(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.JournalEntryCancelledEvent,
	}
}

func (h *atxAccountingEventHandler) OnJournalEntrySubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnJournalEntrySubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusJournalEntryEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.atxAccountingRepo.OnJournalEntrySubmitted(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.JournalEntrySubmittedEvent,
	}
}

func (h *atxAccountingEventHandler) OnPaymentCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnPaymentCancelled"))
				}
			}()
			payload, ok := e.Data.(event.StatusPaymentEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.atxAccountingRepo.OnPaymentCancelled(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.PaymentCancelledEvent,
	}
}

func (h *atxAccountingEventHandler) OnPaymentSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnPaymentSubmitted"))
				}
			}()
			payload, ok := e.Data.(event.StatusPaymentEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.atxAccountingRepo.OnPaymentSubmitted(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		},
		AbortOnError: true,
		Matcher:      domain.PaymentSubmittedEvent,
	}
}
