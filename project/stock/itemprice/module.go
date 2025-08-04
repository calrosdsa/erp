package itemprice

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	itemprice_event "erp/project/stock/itemprice/handler/event"
	itemprice_rest "erp/project/stock/itemprice/handler/rest"
	itemprice_repo "erp/project/stock/itemprice/repository"
	itemprice_ucase "erp/project/stock/itemprice/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {

	itemPriceRepo := itemprice_repo.NewItemPriceService(svc.DBConn(), svc.Helpers())
	itemPriceEventRepo := itemprice_repo.NewItemPriceEventRepo(svc.Helpers(),itemPriceRepo)
	itemPriceUcase := itemprice_ucase.NewItemPriceUseCase(
		svc.Logger(), svc.PermissionService(), svc.CoreService(), itemPriceRepo)

	itemprice_rest.NewItemPriceHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService(), itemPriceUcase,
	)
	itemprice_event.NewEventHandler(svc.Logger(), itemPriceEventRepo, svc.EventBus())

	svc.Container().AddSingleton(domain.ItemPriceUseCase, func(c di.Container) (any, error) {
		return itemPriceUcase, nil
	})
	return nil
}
