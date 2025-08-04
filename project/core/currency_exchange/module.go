package currency_exchange

import (
	"context"
	"erp/pkg/system"
	currency_exchange_rest "erp/project/core/currency_exchange/internal/handler/rest"
	currency_exchange_fsm "erp/project/core/currency_exchange/internal/pkg/fsm"
	currency_exchange_repo "erp/project/core/currency_exchange/internal/repository"
	currency_exchange_ucase "erp/project/core/currency_exchange/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) (err error) {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	fsm := currency_exchange_fsm.NewCurrencyExchangeFsm()
	currencyExchangeRepo := currency_exchange_repo.NewCurrencyExchangeRepo(svc.DBConn(), svc.Helpers())
	currencyExchangeUcase := currency_exchange_ucase.NewCurrencyExchangeUcase(
		svc.Logger(), currencyExchangeRepo, svc.PermissionService(), svc.CoreService(),fsm,
	)
	currency_exchange_rest.NewHandler(svc.HumaApi(),
		svc.Helpers(), currencyExchangeUcase, huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService())
	return nil
}
