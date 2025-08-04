package order

import (
	"context"
	"erp/api/middlewares"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	order_event "erp/project/order/handler/event"
	order_rest "erp/project/order/handler/rest"
	order_fsm "erp/project/order/internal/pkg/fsm"
	order_repo "erp/project/order/repository"
	order_usecase "erp/project/order/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {

	orderEventRepo := order_repo.NewOrderEventRepository(svc.DBConn())
	orderRepo := order_repo.NewOrderRepository(svc.DBConn(), svc.Helpers())

	orderFsm := order_fsm.NewOrderFsm()

	orderUseCase := order_usecase.NewOrderUseCase(
		svc.Helpers(),
		svc.PermissionService(),
		svc.Logger(),
		orderRepo,
		orderFsm,
		svc.CoreService(),
		svc.Container(),
		svc.EventBus(),
		svc.DocumentService(),
	)
	m := middlewares.NewMiddlewares(svc.SessionService(), svc.HumaApi(), svc.Helpers().Jwt)
	order_rest.NewOrderHandler(svc.HumaApi(), svc.Helpers(), huma.Middlewares{m.Authenticate},
		svc.PermissionService(), orderUseCase)

	order_event.NewOrderEventHandler(svc.EventBus(), svc.Logger(), orderEventRepo)

	svc.Container().AddSingleton(domain.OrderUseCase, func(c di.Container) (any, error) {
		return orderUseCase, nil
	})
	return nil
}
