package booking

import (
	"context"
	"erp/pkg/system"
	booking_events "erp/project/regate/booking/internal/handler/event"
	booking_rest "erp/project/regate/booking/internal/handler/rest"
	booking_fsm "erp/project/regate/booking/internal/pkg/fsm"
	booking_repo "erp/project/regate/booking/internal/repository"
	booking_ucase "erp/project/regate/booking/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {	
	bookingFsm := booking_fsm.NewBookingFsm()
	bookingRepo := booking_repo.NewBookingRepository(svc.DBConn(), svc.Helpers(), svc.SettingService(),
	svc.AccountingService())
	bookingRepoSlot := booking_repo.NewBookingSlotRepository(svc.DBConn(), svc.Helpers())
	bookingEventRepo := booking_repo.NewBookingEventRepo(svc.Helpers(), svc.DBConn(), bookingRepoSlot, svc.SettingService(),
svc.Container())
	bookingUseCase := booking_ucase.NewBookingUseCase(
		svc.Logger(), svc.PermissionService(), svc.CoreService(), svc.EventBus(), svc.Container(), bookingRepo,
		bookingRepoSlot, bookingFsm,
	)
	bookingSlotUcase := booking_ucase.NewBookingSlotUseCase(svc.Logger(), bookingRepoSlot)

	booking_rest.NewBookingHandler(
		svc.HumaApi(), svc.Helpers(), svc.PermissionService(),
		huma.Middlewares{svc.Middlewares().Authenticate}, bookingUseCase,
	)
	booking_rest.NewBookingSlotHandler(
		svc.HumaApi(), svc.Helpers(), svc.PermissionService(),
		huma.Middlewares{svc.Middlewares().Authenticate}, bookingSlotUcase,
	)

	booking_events.NewBookingEventHandler(
		svc.EventBus(),
		svc.Logger(),
		bookingEventRepo,
	)
	return nil
}
