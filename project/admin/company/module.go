package admin_company

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	a_company_rest "erp/project/admin/company/internal/handler/rest"
	a_company_repo "erp/project/admin/company/internal/repository"
	a_company_ucase "erp/project/admin/company/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()
	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(),nil
	})
	aCompanyRepo := a_company_repo.NewAdminCompanyRepository(
		svc.DBConn(), svc.Helpers(),svc.Config().PG)
	aCompanyUcase := a_company_ucase.NewAdminCompanyUCase(
		svc.Logger(), aCompanyRepo,svc.EventBus(),container,
	)
	a_company_rest.NewCompanyAHandler(
		svc.HumaApi(), svc.Helpers(), aCompanyUcase, huma.Middlewares{svc.Middlewares().AuthenticateAdmin},
	)
	return nil
}
