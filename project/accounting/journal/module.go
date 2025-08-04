package journal

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	rest_journal "erp/project/accounting/journal/internal/handler/rest"
	journal_entry_fsm "erp/project/accounting/journal/internal/pkg/fsm"
	journal_repo "erp/project/accounting/journal/internal/repository"
	journal_ucase "erp/project/accounting/journal/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()
	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(), nil
	})
	journalEntryFsm := journal_entry_fsm.NewJournalEntryFsm()
	journalRepo := journal_repo.NewJournalRepo(svc.DBConn(), svc.Helpers())
	journalUcase := journal_ucase.NewJournalUseCase(svc.Logger(), svc.PermissionService(), svc.CoreService(),
		journalRepo, container, svc.EventBus(), journalEntryFsm)
	rest_journal.NewJournalHandler(svc.HumaApi(), svc.Helpers(), journalUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())
	return nil
}
