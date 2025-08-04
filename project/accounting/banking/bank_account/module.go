package bank_account

import (
	"context"
	"erp/pkg/system"
	bank_account_rest "erp/project/accounting/banking/bank_account/handler"
	bank_account_fsm "erp/project/accounting/banking/bank_account/internal/pkg/fsm"
	bank_account_repo "erp/project/accounting/banking/bank_account/repository"
	bank_account_usecase "erp/project/accounting/banking/bank_account/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	fsm := bank_account_fsm.NewFsm()
	repo := bank_account_repo.NewRepository(svc.DBConn(), svc.Helpers())
	usecase := bank_account_usecase.NewUseCase(svc.Logger(), svc.CoreService(), svc.PermissionService(), repo,
		fsm)
	bank_account_rest.NewHandler(svc.HumaApi(), svc.Helpers(), huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService(), usecase)
	return nil
}
