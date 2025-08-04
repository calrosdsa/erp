package cost_center

import (
	"context"
	"erp/pkg/system"
	rest_cost_center "erp/project/accounting/cost_center/internal/handler/rest"
	cost_center_repo "erp/project/accounting/cost_center/internal/repository"
	cost_center_ucase "erp/project/accounting/cost_center/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	costCenterRepo := cost_center_repo.NewCostCenterRepo(svc.DBConn(), svc.Helpers())
	costCenterUcase := cost_center_ucase.NewCostCenterUcase(svc.Logger(), costCenterRepo,
		svc.PermissionService(), svc.CoreService())
	rest_cost_center.NewCostCenterHandler(svc.HumaApi(), svc.Helpers(), costCenterUcase,
		huma.Middlewares{svc.Middlewares().Authenticate},svc.PermissionService())
	return nil
}

