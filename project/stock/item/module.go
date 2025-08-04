package item

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	item_rest "erp/project/stock/item/handler/rest"
	item_fsm "erp/project/stock/item/pkg/fsm"
	item_repo "erp/project/stock/item/repository"
	item_ucase "erp/project/stock/item/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {

	// container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
	// 	return svc.DBConn().GetQ().Begin(),nil
	// })

	itemFsm := item_fsm.NewFsm()

	itemRepo := item_repo.NewItemRepository(
		svc.DBConn(), svc.Helpers(),
	)
	itemUcase := item_ucase.NewItemUseCase(
		svc.Logger(), svc.PermissionService(), svc.CoreService(),
		itemRepo, svc.Container(), svc.EventBus(),itemFsm,
	)
	item_rest.NewItemHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, itemUcase, svc.PermissionService())

	svc.Container().AddSingleton(domain.ItemUseCase, func(c di.Container) (any, error) {
		return itemUcase, nil
	})
	return nil
}
