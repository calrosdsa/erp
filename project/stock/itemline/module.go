package itemline

import (
	"context"
	"erp/api/middlewares"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	itemline_event "erp/project/stock/itemline/internal/handler/event"
	itemline_rest "erp/project/stock/itemline/internal/handler/rest"
	itemline_repo "erp/project/stock/itemline/internal/repository"
	itemline_ucase "erp/project/stock/itemline/internal/usecase"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module)Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	fmt.Println("STARTING ITEMLINE MODULE...")
	container := di.New()

	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(),nil
	})
	itemLineRepo := itemline_repo.NewItemLineRepository(svc.DBConn(),svc.Helpers())
	itemLineEventRepo := itemline_repo.NewItemLineEventRepo(svc.DBConn(),svc.Helpers(),itemLineRepo)
	itemLineUseCase := itemline_ucase.NewItemLineUseCase(
		svc.Logger(), svc.PermissionService(), itemLineRepo,svc.EventBus(),container,
	)
	m := middlewares.NewMiddlewares(svc.SessionService(), svc.HumaApi(), svc.Helpers().Jwt)
	itemline_rest.NewItemLineHandler(
		svc.HumaApi(),svc.Helpers(),huma.Middlewares{m.Authenticate},svc.PermissionService(),itemLineUseCase,
	)

	itemline_event.NewItemLineEventHandler(
		svc.EventBus(),svc.Logger(),itemLineEventRepo,
	)
	return nil
}
