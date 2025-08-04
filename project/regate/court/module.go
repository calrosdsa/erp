package court

import (
	"context"
	"erp/pkg/system"
	court_rest "erp/project/regate/court/internal/handler/rest"
	court_repo "erp/project/regate/court/internal/repository"
	court_ucase "erp/project/regate/court/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context,svc system.Service) error {
	return Root(ctx,svc)
}

func Root(ctx context.Context,svc system.Service)error{
	courtRepo := court_repo.NewCourtRepositository(svc.DBConn(),svc.Helpers())
	courtRateRepo := court_repo.NewCourtRateRepository(svc.DBConn(),svc.Helpers())
	courtUseCase := court_ucase.NewCourtUseCase(
		svc.Logger(),svc.PermissionService(),courtRepo,svc.CoreService())
	courtRateUseCase := court_ucase.NewCourtRateUseCase(
		svc.PermissionService(),svc.Logger(),courtRateRepo,
	)
	court_rest.NewCourtHandler(
		svc.HumaApi(),svc.Helpers(),svc.PermissionService(),
		huma.Middlewares{svc.Middlewares().Authenticate},courtUseCase,
	)
	court_rest.NewCourtRateHandler(
		svc.HumaApi(),svc.Helpers(),svc.PermissionService(),
		huma.Middlewares{svc.Middlewares().Authenticate},courtRateUseCase,
	)
	return nil
}