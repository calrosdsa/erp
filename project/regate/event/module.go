package event

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	eventbooking_rest "erp/project/regate/event/internal/handler/rest"
	event_fsm "erp/project/regate/event/internal/pkg/fsm"
	event_repo "erp/project/regate/event/internal/repository"
	event_ucase "erp/project/regate/event/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()

	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(), nil
	})
	fsm := event_fsm.NewEventFsm()
	eventBookingRepo := event_repo.NewEventBookingRepository(svc.DBConn(), svc.Helpers())
	eventBookingUseCase := event_ucase.NewEventBookingUseCase(
		svc.Logger(), svc.PermissionService(), svc.CoreService(), eventBookingRepo,
		fsm, container, svc.EventBus(),
	)
	eventbooking_rest.NewEventBookingHandler(
		svc.HumaApi(), svc.Helpers(), svc.PermissionService(),
		huma.Middlewares{svc.Middlewares().Authenticate}, eventBookingUseCase,
	)
	return nil
}
