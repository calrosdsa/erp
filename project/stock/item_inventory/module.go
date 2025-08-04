package item_inventory

import (
	"context"
	"erp/pkg/system"
	item_inventory_event "erp/project/stock/item_inventory/internal/handler/event"
	item_inventory_rest "erp/project/stock/item_inventory/internal/handler/rest"
	item_inventory_repo "erp/project/stock/item_inventory/internal/repository"
	item_inventory_ucase "erp/project/stock/item_inventory/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) (err error) {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	itemInventoryRepo := item_inventory_repo.NewItemInventoryRepo(
		svc.DBConn(), svc.Helpers(),
	)
	itemInventoryUcase := item_inventory_ucase.NewItemInventoryUcase(
		svc.Logger(), svc.PermissionService(), itemInventoryRepo,
	)
	item_inventory_rest.NewHandler(svc.HumaApi(), svc.Helpers(), itemInventoryUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())
	item_inventory_event.NewItemInventoryEvent(svc.Logger(),itemInventoryRepo,svc.EventBus())
	
	return nil
}
