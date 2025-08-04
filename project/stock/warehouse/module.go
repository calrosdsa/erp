package warehouse

import (
	"context"
	"erp/pkg/system"
	warehouse_rest "erp/project/stock/warehouse/internal/handler/rest"
	warehouse_repo "erp/project/stock/warehouse/internal/repository"
	warehouse_ucase "erp/project/stock/warehouse/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	warehouseRepo := warehouse_repo.NewWareHouseReposotory(
		svc.DBConn(), svc.Helpers(),
	)
	warehouseUseCase := warehouse_ucase.NewWarehouseUseCase(
		svc.Logger(), svc.PermissionService(), svc.CoreService(),
		warehouseRepo,
	)
	warehouse_rest.NewWareHouseHandler(
		svc.HumaApi(), svc.Helpers(), huma.Middlewares{svc.Middlewares().Authenticate},
		warehouseUseCase, svc.PermissionService(),
	)
	return nil
}
