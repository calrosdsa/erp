package account_module

import (
	"context"
	"erp/pkg/system"
	rest_account "erp/project/auth/account/internal/handler/rest"
	account_repo "erp/project/auth/account/internal/repository"
	account_ucase "erp/project/auth/account/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct {
}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	accountRepo := account_repo.NewAccountRepository(
		svc.DBConn(), svc.Helpers(),svc.Config(),
	)
	accountUcase := account_ucase.NewAccountUseCase(
		svc.Logger(), accountRepo,svc.EventBus(),svc.SessionService(),
	)
	rest_account.NewAccountHandler(svc.HumaApi(), svc.Helpers(), svc.SessionService(),
		huma.Middlewares{svc.Middlewares().Authenticate}, accountUcase)
	return nil
}
