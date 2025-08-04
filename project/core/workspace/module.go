package workspace

import (
	"context"
	"erp/pkg/system"
	workspace_rest "erp/project/core/workspace/handler/rest"
	workspace_fsm "erp/project/core/workspace/pkg/fsm"
	workspace_repo "erp/project/core/workspace/repository"
	workspace_ucase "erp/project/core/workspace/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	fsm := workspace_fsm.NewFsm()
	workspaceRepo := workspace_repo.NewWorkSpaceRepo(
		svc.DBConn(), svc.Helpers(),
	)
	workSpaceUcase := workspace_ucase.NewWorkSpaceUseCase(
		svc.PermissionService(), svc.CoreService(), workspaceRepo, svc.Logger(),
		svc.Container(), fsm,
	)
	workspace_rest.NewHandler(svc.HumaApi(), svc.Helpers(), huma.Middlewares{
		svc.Middlewares().Authenticate,
	}, svc.PermissionService(), workSpaceUcase)

	return nil
}
