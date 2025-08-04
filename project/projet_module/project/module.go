package project

import (
	"context"
	"erp/pkg/system"
	rest_project "erp/project/projet_module/project/internal/handler/rest"
	project_repo "erp/project/projet_module/project/internal/repository"
	project_ucase "erp/project/projet_module/project/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	projectRepo := project_repo.NewProjectRepository(svc.DBConn(), svc.Helpers())
	projectUcase := project_ucase.NewProjectUcase(svc.Logger(), projectRepo,
		svc.PermissionService(), svc.CoreService())
	rest_project.NewProjectHandler(svc.HumaApi(), svc.Helpers(), projectUcase,
		huma.Middlewares{svc.Middlewares().Authenticate},svc.PermissionService())
	return nil
}

