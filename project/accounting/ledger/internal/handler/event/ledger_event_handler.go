package ledger_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	ledger_repo "erp/project/accounting/ledger/internal/repository"
)

type ledgerEventHandler struct {
	emitLog logger.EmitLog
	ledgerEventRepo ledger_repo.LedgerEventRepo
}

func NewLedgerEventHandler(
	bus bus.Bus,
	logger  logger.Logger,
	ledgerEventRepo ledger_repo.LedgerEventRepo,
){
	h := ledgerEventHandler{
		emitLog: logger.EmitLog("ledger-event-handler"),
		ledgerEventRepo: ledgerEventRepo,
	}
	bus.RegisterHandler(domain.EventCompanyCreated, h.OnCompanyCreated())
}

func (h *ledgerEventHandler) OnCompanyCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.CreatedCompanyEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.ledgerEventRepo.CreateChartOfAccountsCompany(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		},
		Matcher: domain.EventCompanyCreated,
		AbortOnError: true,
	}
}
