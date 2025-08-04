package ledger

import (
	"context"
	"erp/pkg/system"
	ledger_event "erp/project/accounting/ledger/internal/handler/event"
	ledger_rest "erp/project/accounting/ledger/internal/handler/rest"
	pkg_ledger "erp/project/accounting/ledger/internal/pkg"
	ledger_repo "erp/project/accounting/ledger/internal/repository"
	ledger_ucase "erp/project/accounting/ledger/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	accountBuilder := pkg_ledger.NewAccountBuilder(svc.Helpers().Locale)
	ledgerEventRepo := ledger_repo.NewLedgerEventRepo(svc.Helpers(), accountBuilder)
	ledgerRepo := ledger_repo.NewLedgerRepository(
		svc.DBConn(), svc.Helpers(),
	)
	ledgerUseCase := ledger_ucase.NewLedgerUseCase(
		ledgerRepo, svc.Logger(), svc.PermissionService(),svc.CoreService(),
	)
	ledger_rest.NewLedgerHandler(
		svc.HumaApi(), huma.Middlewares{svc.Middlewares().Authenticate}, svc.Helpers(),
		ledgerUseCase, svc.PermissionService(),
	)
	ledger_event.NewLedgerEventHandler(svc.EventBus(), svc.Logger(), ledgerEventRepo)

	return nil
}
