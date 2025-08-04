package company

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	rest_company "erp/project/company/internal/handler/rest"
	company_repo "erp/project/company/internal/repository"
	company_ucase "erp/project/company/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()

	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(), nil
	})
	companyRepo := company_repo.NewCompanyRepository(svc.DBConn(), svc.Helpers(), svc.Logger())
	companyUCase := company_ucase.NewCompanyUseCase(
		svc.Logger(), companyRepo, svc.PermissionService(), svc.CoreService(),
		svc.EventBus(), container,
	)
	rest_company.NewCompanyHandler(
		svc.HumaApi(), svc.Helpers(), svc.PermissionService(), companyUCase,
		huma.Middlewares{svc.Middlewares().Authenticate},
	)

	return nil
}
