package bank

import (
	"context"
	"erp/pkg/system"
	bank_rest "erp/project/accounting/banking/bank/handler"
	bank_fsm "erp/project/accounting/banking/bank/internal/pkg/fsm"
	bank_repo "erp/project/accounting/banking/bank/repository"
	bank_usecase "erp/project/accounting/banking/bank/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	fsm := bank_fsm.NewFsm()
	repo := bank_repo.NewRepository(svc.DBConn(), svc.Helpers())
	usecase := bank_usecase.NewUseCase(
		svc.Logger(), svc.CoreService(), svc.PermissionService(), repo, fsm,
	)
	bank_rest.NewHandler(
		svc.HumaApi(), svc.Helpers(), huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService(), usecase,
	)
	return nil
}
