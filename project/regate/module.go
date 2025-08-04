package regate

import (
	"context"
	"erp/pkg/system"
	"erp/project/regate/booking"
	"erp/project/regate/chart"
	"erp/project/regate/court"
	"erp/project/regate/event"
	regate_domain "erp/project/regate/internal/domain"
)

type Project struct{}

type monolith struct {
	system.Service
	modules []system.Module
}

func (m Project) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}
func Root(ctx context.Context, svc system.Service) error {

	//Register topics event
	svc.EventBus().RegisterTopics(
		regate_domain.BookingCompletedEvent,
		regate_domain.BookingCancelEvent,
		regate_domain.EditPaidBookingEvent,
		regate_domain.BookingRescheduleEvent,
		regate_domain.BookingDeletedEvent,
		regate_domain.CancelEventBooking,
		regate_domain.CompletedEventBooking,
	)

	m := monolith{
		Service: svc,
		modules: []system.Module{
			&court.Module{},
			&booking.Module{},
			&event.Module{},
			&chart.Module{},
		},
	}

	return m.startupModules()
}

func (m *monolith) startupModules() error {

	for _, module := range m.modules {
		ctx := m.Waiter().Context()
		if err := module.Startup(ctx, m); err != nil {
			return err
		}
	}
	return nil
}
