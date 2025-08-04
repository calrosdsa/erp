package booking_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	contact_ucase "erp/project/core/contact/usecase"
	booking_repo "erp/project/regate/booking/internal/repository"
	regate_domain "erp/project/regate/internal/domain"
	regate_event "erp/project/regate/internal/domain/event"
	"fmt"
	"strconv"
)

type BookingUseCase interface {
	CreateBooking(req *common.RequestContext, d dto.CreateBookingBody) error
	GetBookings(req *common.RequestContext, i *dto.RequestBookings) (
		res dto.PaginationResult[[]dto.BookingDto], err error)
	GetBooking(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.BookingDto], err error)

	ValidateBooking(req *common.RequestContext, i dto.ValidateBookingData) (
		res dto.ValidateBookingData, err error)

	UpdateBookingState(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
	EditBookingPaidAmount(req *common.RequestContext, d dto.BookingPaymentBody) error
	BookingReschedule(req *common.RequestContext, d dto.BookingRescheduleBody) (err error)
	UpdateBookingBatch(req *common.RequestContext, d *dto.UpdateBookingBatchRequest) (
		err error)
}

type bookingUseCase struct {
	emitLog         logger.EmitLog
	bookingRepo     booking_repo.BookingRepository
	bookingSlotRepo booking_repo.BookingSlotRepository
	permission      repository.PermissionService
	core            repository.CoreService
	c               di.Container
	bus             bus.Bus
	fsm             fsm.FsmState
	contactUseCase  contact_ucase.ContactUseCase
}

func NewBookingUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	core repository.CoreService,
	bus bus.Bus,
	c di.Container,
	bookingRepo booking_repo.BookingRepository,
	bookingSlotRepo booking_repo.BookingSlotRepository,
	fsm fsm.FsmState,
) BookingUseCase {
	return &bookingUseCase{
		emitLog:         logger.EmitLog("booking-usecase"),
		permission:      permission,
		bookingRepo:     bookingRepo,
		bookingSlotRepo: bookingSlotRepo,
		bus:             bus,
		c:               c,
		core:            core,
		fsm:             fsm,
		contactUseCase:  c.Get(domain.ContactUseCase).(contact_ucase.ContactUseCase),
	}
}

func (u *bookingUseCase) UpdateBookingBatch(req *common.RequestContext, d *dto.UpdateBookingBatchRequest) (
	err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateBookingBatch"))
		}
		err = u.closeTx(tx, err)
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, regate_domain.BOOKING, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	bookings, err := u.bookingRepo.UpdateBookingBatch(req, d)
	if err != nil {
		return
	}
	payload := regate_event.BookingStatusData{
		Profile:         req.Profile,
		Bookings:        bookings,
		Tx:              tx,
		CompanyDefaults: req.CompanyDefaults,
	}
	switch d.Body.TargetState {
	case proto.State_CANCELLED.String():
		err = u.bus.Emit(req.Ctx, regate_domain.BookingCancelEvent, payload)
	case proto.State_COMPLETED.String():
		err = u.bus.Emit(req.Ctx, regate_domain.BookingCompletedEvent, payload)
	case proto.State_DELETED.String():
		err = u.bus.Emit(req.Ctx, regate_domain.BookingDeletedEvent, payload)
	}
	return
}

func (u *bookingUseCase) BookingReschedule(req *common.RequestContext, d dto.BookingRescheduleBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("BookingReschedule"))
		}
		err = u.closeTx(tx, err)

	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, regate_domain.BOOKING, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.bookingRepo.BookingReschedule(tx, req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, regate_domain.BookingRescheduleEvent, regate_event.RescheduleBookingEventData{
		Tx:                tx,
		Company:           req.ActiveCompany,
		BookingReschedule: d,
	})
	return
}

func (u *bookingUseCase) EditBookingPaidAmount(req *common.RequestContext, d dto.BookingPaymentBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditBookingPaidAmount"))
		}
		err = u.closeTx(tx, err)

	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, regate_domain.BOOKING, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	booking, err := u.bookingRepo.EditBookingPaidAmount(tx, req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, regate_domain.EditPaidBookingEvent, regate_event.EditPaidBookingEventData{
		Booking: booking,
		Tx:      tx,
		Profile: req.Profile,
	})
	return
}

func (u *bookingUseCase) ValidateBooking(req *common.RequestContext, i dto.ValidateBookingData) (
	res dto.ValidateBookingData, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ValidateBooking"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, regate_domain.BOOKING, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.bookingRepo.ValidateBooking(req, i)
	return
}

func (u *bookingUseCase) CreateBooking(req *common.RequestContext, d dto.CreateBookingBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateBooking"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	if allow := u.permission.CheckPermission(req.Ctx, req, regate_domain.BOOKING, domain.CREATE); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.bookingRepo.CreateBooking(tx, req, d)
	if err != nil {
		return
	}
	err = u.bookingSlotRepo.UpdatedCreatedBookings(tx, req.Ctx, req.ActiveCompany.ID)
	// res = dto.BookingDtoFromModel(&booking)
	return
}
func (r *bookingUseCase) GetBookings(req *common.RequestContext, i *dto.RequestBookings) (
	res dto.PaginationResult[[]dto.BookingDto], err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetBookings"))
		}
	}()
	if allow := r.permission.CheckPermission(req.Ctx, req, regate_domain.BOOKING, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = r.bookingRepo.GetBookings(req, i)
	return
}

func (r *bookingUseCase) GetBooking(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.BookingDto], err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetBooking"))
		}
	}()
	if allow := r.permission.CheckPermission(req.Ctx, req, regate_domain.BOOKING, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = r.bookingRepo.GetBooking(req, i)
	if err != nil {
		return
	}
	//Getting activities
	res.Activities = r.core.GerActivitiesByPartyID(req, res.Entity.ID)
	contacts, err := r.contactUseCase.GetContacts(req, dto.ContactsRequest{
		PartyID: strconv.Itoa(int(res.Entity.PartyID)),
	})
	if err != nil {
		return
	}
	fmt.Println("CONTACTS", contacts)
	res.Contacts = contacts
	return
}

func (u *bookingUseCase) UpdateBookingState(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateBookingState"))
		}
		err = u.closeTx(tx, err)

	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, regate_domain.BOOKING, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := u.fsm.NextState(i.Body.CurrentState, i.Body.Events)
	if err != nil {
		return err
	}
	booking, err := u.bookingRepo.UpdateBookingState(tx, req, i.Body.PartyID, i.Body.CurrentState, nextState)
	if err != nil {
		return
	}
	payload := regate_event.BookingStatusData{
		Profile:         req.Profile,
		Bookings:        []*model.Booking{booking},
		Tx:              tx,
		CompanyDefaults: req.CompanyDefaults,
	}
	switch nextState {
	case proto.State_CANCELLED.String():
		fmt.Println("CANCELLED STATE", booking.ID, tx)
		err = u.bus.Emit(req.Ctx, regate_domain.BookingCancelEvent, payload)
	case proto.State_COMPLETED.String():
		err = u.bus.Emit(req.Ctx, regate_domain.BookingCompletedEvent, payload)
	case proto.State_DELETED.String():
		err = u.bus.Emit(req.Ctx, regate_domain.BookingDeletedEvent, payload)
	}
	return
}

func (s *bookingUseCase) closeTx(tx *query.QueryTx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
