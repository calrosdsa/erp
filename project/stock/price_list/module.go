package pricelist

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	price_list_event "erp/project/stock/price_list/handler/event"
	rest_price_list "erp/project/stock/price_list/handler/rest"
	price_list_repo "erp/project/stock/price_list/repository"
	price_list_ucase "erp/project/stock/price_list/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}
func Root(ctx context.Context, svc system.Service) error {
	priceListEventRepo := price_list_repo.NewPriceListEventRepo()
	priceListRepo := price_list_repo.NewPriceListServer(svc.DBConn(), svc.Helpers(), svc.Logger())
	priceListUcase := price_list_ucase.NewPriceListUseCase(svc.Logger(), priceListRepo,
		svc.PermissionService(), svc.CoreService())
	rest_price_list.NewPriceListHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, priceListUcase, svc.PermissionService())
	price_list_event.NewPriceListEventHandler(svc.Logger(),priceListEventRepo,svc.EventBus())

	svc.Container().AddSingleton(domain.PriceListUseCase, func(c di.Container) (any, error) {
		return priceListUcase, nil
	})
	return nil
}
