package supplier

import (
	"context"
	// "erp/internal/domain"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	supplier_rest "erp/project/buying/supplier/internal/handler/rest"
	supplier_fsm "erp/project/buying/supplier/internal/pkg/fsm"
	supplier_repo "erp/project/buying/supplier/internal/repository"
	supplier_ucase "erp/project/buying/supplier/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()
	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(),nil
	})
	fsm := supplier_fsm.NewSupplierFsm()
	supplierRepo := supplier_repo.NewSupplierRepository(svc.DBConn(), svc.Helpers())
	supplierUseCase := supplier_ucase.NewSupplierUseCase(
		svc.PermissionService(), supplierRepo, svc.Logger(), svc.CoreService(),
		svc.EventBus(),container,fsm,
	)
	supplier_rest.NewSupplierHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, supplierUseCase,
		svc.PermissionService())
	return nil
}
