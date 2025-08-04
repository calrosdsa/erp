package deal

import (
	"context"
	"erp/pkg/system"
	deal_mcp "erp/project/crm/deal/handler/mcp"
	deal_rest "erp/project/crm/deal/handler/rest"
	deal_repo "erp/project/crm/deal/repository"
	deal_ucase "erp/project/crm/deal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	dealRepo := deal_repo.NewDealRepository(svc.DBConn(), svc.Helpers())
	dealUcase := deal_ucase.NewDealUseCase(
		svc.PermissionService(), svc.CoreService(), dealRepo, svc.Logger(), svc.EventBus(),
		svc.Container(),
	)
	deal_rest.NewHandler(
		svc.HumaApi(), svc.Helpers(), huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService(), dealUcase,
	)
	deal_mcp.NewDealMcp(svc.Mcp(), svc.Helpers(), svc.PermissionService(), dealUcase)
	return nil
}
