package event_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	event_repo "erp/project/regate/event/internal/repository"
	regate_domain "erp/project/regate/internal/domain"
	regate_event "erp/project/regate/internal/domain/event"
	"fmt"
)

type EventBookingUseCase interface {
	CreateEventBooking(req *common.RequestContext, i dto.EventBookingData) (res dto.EventBookingDto, err error)
	GetEventBookings(req *common.RequestContext, i *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.EventBookingDto], err error)
	GetEventBooking(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.EventBookingDetail], err error)
	EditEvent(req *common.RequestContext, d dto.EventBookingData) (err error)
	UpdateStatus(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error)
	DeleteEventBatch(req *common.RequestContext, i *dto.DeleteEventBatchRequest) (err error)
}

type eventBookingUseCase struct {
	emitLog          logger.EmitLog
	permission       repository.PermissionService
	fsm              fsm.FsmState
	core             repository.CoreService
	eventBookingRepo event_repo.EventBookingRepository
	c                di.Container
	bus              bus.Bus
}

func NewEventBookingUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	core repository.CoreService,
	eventBookingRepo event_repo.EventBookingRepository,
	fsm fsm.FsmState,
	c di.Container,
	bus bus.Bus,
) EventBookingUseCase {
	return &eventBookingUseCase{
		emitLog:          logger.EmitLog("event-booking-usecase"),
		permission:       permission,
		core:             core,
		eventBookingRepo: eventBookingRepo,
		fsm:              fsm,
		c:                c,
		bus:              bus,
	}
}

func (r *eventBookingUseCase) DeleteEventBatch(req *common.RequestContext, i *dto.DeleteEventBatchRequest) (err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("DeleteEventBatch"))
		}
	}()
	err = r.permission.CheckPermissionEntity(req.Ctx, req, regate_domain.EVENT_BOOKING, domain.DELETE)
	if err != nil {
		return domain.ACTION_NOT_ALLOWED
	}
	err = r.eventBookingRepo.DeleteEventBatch(req, i)
	return
}

func (r *eventBookingUseCase) UpdateStatus(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error) {
	ctx := r.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("UpstateEventState"))
		}
		err = r.closeTx(tx, err)
	}()
	err = r.permission.CheckPermissionEntity(req.Ctx, req, regate_domain.EVENT_BOOKING, domain.EDIT)
	if err != nil {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := r.fsm.NextState(i.Body.CurrentState, i.Body.Events, i.Body.PartyType)
	if err != nil {
		return err
	}
	fmt.Println("NEXT STATE", nextState)
	event, err := r.eventBookingRepo.UpdateStatus(req, tx, i.Body.PartyID, i.Body.CurrentState, nextState)
	if err != nil {
		return
	}
	payload := regate_event.StatusEventBookingData{
		Event:           event,
		Tx:              tx,
		Profile:         req.Profile,
		CompanyDefaults: req.CompanyDefaults,
	}
	switch nextState {
	case proto.State_CANCELLED.String():
		err = r.bus.Emit(req.Ctx, regate_domain.CancelEventBooking, payload)
	case proto.State_COMPLETED.String():
		err = r.bus.Emit(req.Ctx, regate_domain.CompletedEventBooking, payload)
	}
	return
}
func (r *eventBookingUseCase) EditEvent(req *common.RequestContext, d dto.EventBookingData) (err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("EditEvent"))
		}
	}()
	err = r.permission.CheckPermissionEntity(req.Ctx, req, regate_domain.EVENT_BOOKING, domain.EDIT)
	if err != nil {
		return
	}
	err = r.eventBookingRepo.EditEvent(req, d)
	return
}

func (r *eventBookingUseCase) GetEventBookings(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.EventBookingDto], err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetEventBookings"))
		}
	}()
	if allow := r.permission.CheckPermission(req.Ctx, req, regate_domain.EVENT_BOOKING, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = r.eventBookingRepo.GetEventBookings(req, d)
	return
}
func (r *eventBookingUseCase) GetEventBooking(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.EventBookingDetail], err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetEventBooking"))
		}
	}()
	if allow := r.permission.CheckPermission(req.Ctx, req, regate_domain.EVENT_BOOKING, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = r.eventBookingRepo.GetEventBooking(req, i)

	res.Activities = r.core.GerActivitiesByPartyID(req, res.Entity.EventBooking.ID)
	return
}

func (r *eventBookingUseCase) CreateEventBooking(req *common.RequestContext, i dto.EventBookingData) (
	res dto.EventBookingDto, err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateEventBooking"))
		}
	}()
	if allow := r.permission.CheckPermission(req.Ctx, req, regate_domain.EVENT_BOOKING, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = r.eventBookingRepo.CreateEventBooking(req, i)
	return
}

func (r *eventBookingUseCase) closeTx(tx *query.QueryTx, err error) error {
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
