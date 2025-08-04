package customer

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	customer_mcp "erp/project/selling/customer/internal/handler/mcp"
	customer_rest "erp/project/selling/customer/internal/handler/rest"
	customer_fsm "erp/project/selling/customer/internal/pkg/fsm"
	customer_repo "erp/project/selling/customer/internal/repository"
	customer_ucase "erp/project/selling/customer/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx,svc)
}

func Root(ctx context.Context,svc system.Service)error {
	container := di.New()
	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(),nil
	})
	fms := customer_fsm.NewCustomerFsm()
	customerRepo := customer_repo.NewCustomerRepository(
		svc.DBConn(),svc.Helpers(),
	)
	customerUCase := customer_ucase.NewCustomerUseCase(
		svc.Helpers(),svc.Logger(),customerRepo,svc.PermissionService(),
		svc.EventBus(),container,svc.CoreService(),fms,
	)
	customer_rest.NewCustomerHandler(
		svc.HumaApi(),svc.Helpers(),huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService(),customerUCase,
	)
	
	// Initialize MCP tools
	customer_mcp.NewCustomerMcp(
		svc.Mcp(),
		svc.Helpers(),
		svc.PermissionService(),
		customerUCase,
	)
	
	return nil
}
