package booking_events

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/bus"
	"erp/pkg/logger"
	booking_repo "erp/project/regate/booking/internal/repository"
	regate_domain "erp/project/regate/internal/domain"
	regate_event "erp/project/regate/internal/domain/event"
	"fmt"
)

type bookingEventHandler struct {
	bus              bus.Bus
	emitLog          logger.EmitLog
	bookingEventRepo booking_repo.BookingEventRepository
}

func NewBookingEventHandler(
	bus bus.Bus,
	logger logger.Logger,
	bookingEventRepo booking_repo.BookingEventRepository,
) {
	handler := bookingEventHandler{
		emitLog:          logger.EmitLog("booking-event"),
		bookingEventRepo: bookingEventRepo,
	}
	bus.RegisterHandler(regate_domain.BookingCompletedEvent, handler.OnCompletedBookingEvent())
	bus.RegisterHandler(regate_domain.BookingDeletedEvent, handler.OnDeleteBookingEvent())
	bus.RegisterHandler(regate_domain.BookingCancelEvent, handler.OnCancelBookingEvent())
	bus.RegisterHandler(regate_domain.EditPaidBookingEvent, handler.OnEditPaidBookingEvent())
	bus.RegisterHandler(regate_domain.BookingRescheduleEvent, handler.OnBookingRescheduleEvent())
	bus.RegisterHandler(regate_domain.CancelEventBooking, handler.OnEventCancelled())
	bus.RegisterHandler(regate_domain.CompletedEventBooking, handler.OnEventCompleted())
}

func (h *bookingEventHandler) OnEventCompleted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			payload, ok := e.Data.(regate_event.StatusEventBookingData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnEventCompleted"))
				}
			}()
			err = h.bookingEventRepo.OnEventCompleted(ctx, payload)
			return
		},
		Matcher:      regate_domain.CompletedEventBooking,
		AbortOnError: true,
	}
}

func (h *bookingEventHandler) OnEventCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			payload, ok := e.Data.(regate_event.StatusEventBookingData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnEventCancelled"))
				}
			}()
			err = h.bookingEventRepo.OnEventCancelled(ctx, payload)
			return
		},
		Matcher:      regate_domain.CancelEventBooking,
		AbortOnError: true,
	}
}

func (h *bookingEventHandler) OnBookingRescheduleEvent() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnBookingRescheduleEvent"))
				}
			}()
			payload, ok := e.Data.(regate_event.RescheduleBookingEventData)
			fmt.Println("On edit paid booking event", payload)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.bookingEventRepo.OnRescheduleBooking(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		},
		Matcher:      regate_domain.BookingRescheduleEvent,
		AbortOnError: true,
	}
}

func (h *bookingEventHandler) OnEditPaidBookingEvent() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnEditPaidBookingEvent"))
				}
			}()
			payload, ok := e.Data.(regate_event.EditPaidBookingEventData)
			fmt.Println("On edit paid booking event", payload)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.bookingEventRepo.EditBookingPaidAmount(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		},
		Matcher:      regate_domain.EditPaidBookingEvent,
		AbortOnError: true,
	}
}
func (h *bookingEventHandler) OnDeleteBookingEvent() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnDeleteBookingEvent"))
				}
			}()
			payload, ok := e.Data.(regate_event.BookingStatusData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.bookingEventRepo.DeleteBookings(ctx, payload)
			return nil
		},
		Matcher:      regate_domain.BookingDeletedEvent,
		AbortOnError: true,
	}
}

func (h *bookingEventHandler) OnCompletedBookingEvent() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnCompletedBookingEvent"))
				}
			}()
			payload, ok := e.Data.(regate_event.BookingStatusData)
			fmt.Println("PAYLOAD completed", payload)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.bookingEventRepo.CompletedBookings(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		},
		Matcher:      regate_domain.BookingCompletedEvent,
		AbortOnError: true,
	}
}

func (h *bookingEventHandler) OnCancelBookingEvent() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnCompletedBookingEvent"))
				}
			}()
			payload, ok := e.Data.(regate_event.BookingStatusData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.bookingEventRepo.CancelBookings(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		},
		Matcher:      regate_domain.BookingCancelEvent,
		AbortOnError: true,
	}
}
