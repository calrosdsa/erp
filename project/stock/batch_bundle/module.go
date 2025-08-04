package batch_bundle

import (
	"context"
	"erp/pkg/system"
	rest_batch_bundle "erp/project/stock/batch_bundle/internal/handler/rest"
	batch_bundle_repo "erp/project/stock/batch_bundle/internal/repository"
	batch_bundle_ucase "erp/project/stock/batch_bundle/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	batchBundleRepo := batch_bundle_repo.NewBatchBundleRepository(svc.DBConn(), svc.Helpers())
	batchBundleUcase := batch_bundle_ucase.NewBatchBundleUcase(svc.Logger(), batchBundleRepo,
		svc.PermissionService(), svc.CoreService())
	rest_batch_bundle.NewBatchBundleHandler(svc.HumaApi(), svc.Helpers(), batchBundleUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())
	return nil
}
