package auth_admin

import (
	"context"
	"erp/pkg/system"
	auth_admin_rest "erp/project/admin/auth/internal/handler/rest"
	auth_admin_repo "erp/project/admin/auth/internal/repository"
	auth_admin_ucase "erp/project/admin/auth/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}
func Root(ctx context.Context, svc system.Service) error {
	authAdminRepo := auth_admin_repo.NewAdminAuthRepository(svc.DBConn(), svc.Config())
	roleTemplateRepo := auth_admin_repo.NewRoleTemplateRepo(svc.DBConn(), svc.Helpers())

	authAdminUcase := auth_admin_ucase.NewAdminAuthUseCase(svc.Logger(), authAdminRepo)
	roleTemplateUcase := auth_admin_ucase.NewRoleTemplateUcase(svc.Logger(), roleTemplateRepo)
	auth_admin_rest.NewAuthAdminHandler(
		svc.HumaApi(), svc.Helpers(), authAdminUcase,
	)
	auth_admin_rest.NewRoleTemplateHandler(
		svc.HumaApi(), svc.Helpers(), roleTemplateUcase, huma.Middlewares{svc.Middlewares().AuthenticateAdmin},
	)
	return nil
}
