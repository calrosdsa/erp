package task

import (
	"context"
	"erp/pkg/system"
	task_mcp "erp/project/projet_module/task/handler/mcp"
	task_rest "erp/project/projet_module/task/handler/rest"
	task_repo "erp/project/projet_module/task/repository"
	task_ucase "erp/project/projet_module/task/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	taskRepo := task_repo.NewTaskRepository(svc.DBConn(), svc.Helpers())
	taskUcase := task_ucase.NewTaskUseCase(
		svc.PermissionService(), svc.CoreService(), taskRepo, svc.Logger(), svc.EventBus(),
		svc.Container(),
	)
	task_rest.NewHandler(
		svc.HumaApi(), svc.Helpers(), huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService(), taskUcase,
	)
	task_mcp.NewTaskMcp(svc.Mcp(), svc.Helpers(), svc.PermissionService(), taskUcase)
	return nil
}
