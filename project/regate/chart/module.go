package chart

import (
	"context"
	"erp/pkg/system"
	chart_rest "erp/project/regate/chart/internal/handler/rest"
	chart_repo "erp/project/regate/chart/internal/repository"
	chart_ucase "erp/project/regate/chart/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context,svc system.Service)(error){
	return Root(ctx,svc)
}

func Root(ctx context.Context,svc system.Service)(error){
	chartRepo := chart_repo.NewChartRepository(svc.DBConn())
	chartUseCase := chart_ucase.NewChartUseCase(
		svc.Logger(),chartRepo,
	)
	chart_rest.NewChartHandler(
		svc.HumaApi(),svc.Helpers(),svc.PermissionService(),
		huma.Middlewares{svc.Middlewares().Authenticate},chartUseCase,
	)
	return nil
}