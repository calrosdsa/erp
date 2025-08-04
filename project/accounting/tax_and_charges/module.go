package tax_and_charges

import (
	"context"
	"erp/pkg/system"
	tac_event "erp/project/accounting/tax_and_charges/handler/event"
	tac_rest "erp/project/accounting/tax_and_charges/handler/rest"
	tac_repo "erp/project/accounting/tax_and_charges/repository"
	tac_usecase "erp/project/accounting/tax_and_charges/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	tacRepo := tac_repo.NewTacRepository(svc.DBConn(),svc.Helpers())
	tacEventRepo := tac_repo.NewTacEventRepository(svc.Helpers())
	tacUcase := tac_usecase.NewTacUseCase(svc.Logger(), tacRepo)
	tac_event.NewTacEventHandler(svc.EventBus(), tacEventRepo, svc.Logger())
	tac_rest.NewTacHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService(),
		tacUcase)
	return nil
}
