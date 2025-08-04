package pricing

import (
	"context"
	"erp/pkg/system"
	pricing_rest "erp/project/pricing_module/pricing/internal/handler/rest"
	pricing_fsm "erp/project/pricing_module/pricing/internal/pkg/fsm"
	pricing_repo "erp/project/pricing_module/pricing/internal/repository"
	pricing_ucase "erp/project/pricing_module/pricing/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) (err error) {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) (err error) {

	fsm := pricing_fsm.NewPricingFsm()
	pricingGeneratorRepo := pricing_repo.NewPricingGeneratorRepo(svc.DBConn())
	pricingRepo := pricing_repo.NewPricingRepository(svc.DBConn(), svc.Helpers())
	pricingUcase := pricing_ucase.NewPricingUcase(
		svc.Logger(), pricingRepo, svc.PermissionService(), svc.CoreService(), fsm,
	)
	pricingGeneratorUcase := pricing_ucase.NewPricingGeneratorUcase(
		svc.Container(), svc.Logger(), pricingGeneratorRepo,
		svc.Helpers(),
	)
	pricing_rest.NewHandler(svc.HumaApi(),
		svc.Helpers(), pricingUcase, huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService(), pricingGeneratorUcase)

	return nil
}
