package admin_core

import (
	"context"
	"erp/pkg/system"
	rest_core "erp/project/admin/core/internal/handler/rest"
	pg_core "erp/project/admin/core/internal/repository"
	core_ucase "erp/project/admin/core/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context,svc system.Service) error {
	return Root(ctx,svc)
}

func Root(ctx context.Context,svc system.Service) error {
	entityRepo := pg_core.NewCoreRepository(svc.DBConn(),svc.Helpers())
	entityUCase := core_ucase.NewCoreEntityUseCase(svc.Logger(),entityRepo)
	rest_core.NewCoreEntityHandler(
		svc.HumaApi(),svc.Helpers(),entityUCase,huma.Middlewares{svc.Middlewares().AuthenticateAdmin},
	)
	return nil
}