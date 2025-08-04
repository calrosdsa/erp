package stage

import (
	"context"
	"erp/pkg/system"
	stage_mcp "erp/project/core/stage/handler/mcp"
	stage_rest "erp/project/core/stage/handler/rest"
	stage_repo "erp/project/core/stage/repository"
	stage_ucase "erp/project/core/stage/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	stageRepo := stage_repo.NewStageRepository(svc.DBConn(), svc.Helpers())
	stageUcase := stage_ucase.NewStageUseCase(
		svc.PermissionService(), svc.CoreService(), stageRepo, svc.Logger(), svc.EventBus(),
	)
	stage_rest.NewHandler(
		svc.HumaApi(), svc.Helpers(), huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService(), stageUcase,
	)
	stage_mcp.NewStageMcp(svc.Mcp(), svc.Helpers(), svc.PermissionService(), stageUcase)
	return nil
}
