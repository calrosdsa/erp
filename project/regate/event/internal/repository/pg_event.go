package event_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type EventBookingRepository interface {
	GetEventBookings(req *common.RequestContext, i *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.EventBookingDto], err error)
	GetEventBooking(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.EventBookingDetail], err error)
	CreateEventBooking(req *common.RequestContext, i dto.EventBookingData) (dto.EventBookingDto, error)
	EditEvent(req *common.RequestContext, d dto.EventBookingData) (err error)
	UpdateStatus(req *common.RequestContext, tx *query.QueryTx,
		id, prevState, nextState string) (*model.EventBooking, error)
	DeleteEventBatch(req *common.RequestContext, d *dto.DeleteEventBatchRequest) (err error)
}

type eventBookingRepository struct {
	Q         *query.Query
	conn      db.Connection
	convertor helpers.ConvertorHelper
}

func NewEventBookingRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) EventBookingRepository {
	return &eventBookingRepository{
		conn:      conn,
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
	}
}
func (r *eventBookingRepository) DeleteEventBatch(req *common.RequestContext, d *dto.DeleteEventBatchRequest) (err error) {
	_, err = r.Q.WithContext(req.Ctx).EventBooking.Where(r.Q.EventBooking.ID.In(d.Body.EventIds...)).Delete()
	return
}

func (r *eventBookingRepository) UpdateStatus(req *common.RequestContext, tx *query.QueryTx,
	id, prevState, nextState string) (res *model.EventBooking, err error) {
	eventBookingQ := tx.EventBooking
	_, err = tx.EventBooking.WithContext(req.Ctx).Where(
		eventBookingQ.CompanyID.Eq(req.ActiveCompany.ID),
		eventBookingQ.Status.Eq(prevState),
		eventBookingQ.UUID.Eq(id),
	).UpdateSimple(eventBookingQ.Status.Value(nextState))
	if err != nil {
		return
	}
	res, err = tx.EventBooking.WithContext(req.Ctx).Where(
		eventBookingQ.CompanyID.Eq(req.ActiveCompany.ID),
		eventBookingQ.UUID.Eq(id),
	).First()
	if err != nil {
		return
	}
	return
}

func (r *eventBookingRepository) EditEvent(req *common.RequestContext, d dto.EventBookingData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.EventBooking.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.EventBooking{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.EventBooking.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()

	return
}

func (r *eventBookingRepository) GetEventBookings(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.EventBookingDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	eventBookingQ := r.Q.EventBooking
	builder := r.Q.WithContext(req.Ctx).EventBooking

	//ADDING CONDITIONS
	conds = append(conds, eventBookingQ.CompanyID.Eq(req.ActiveCompany.ID))
	if d.Status != "" {
		conds = append(conds, eventBookingQ.Status.Eq(d.Status))
	}
	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.Payment.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}

	builder = builder.Select(
		eventBookingQ.ID, eventBookingQ.UUID, eventBookingQ.Name, eventBookingQ.CreatedAt,
		eventBookingQ.Status, eventBookingQ.Description,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
func (r *eventBookingRepository) GetEventBooking(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.EventBookingDetail], err error) {
	id := r.convertor.StrtoInt(i.ID)
	eventBookingQ := r.Q.EventBooking
	err = r.Q.WithContext(req.Ctx).EventBooking.Select(
		eventBookingQ.ID, eventBookingQ.UUID, eventBookingQ.Name, eventBookingQ.CreatedAt,
		eventBookingQ.Status, eventBookingQ.Description,
	).Where(
		eventBookingQ.CompanyID.Eq(req.ActiveCompany.ID),
		eventBookingQ.ID.Eq(id),
	).Scan(&res.Entity.EventBooking)
	if err != nil {
		return res, err
	}
	res.Entity.BookingInfo, err = r.getEventBookingInfo(req.Ctx, res.Entity.EventBooking.ID)

	return
}

func (r *eventBookingRepository) getEventBookingInfo(ctx context.Context, eventBookingID int64) (
	res dto.EventBookingInfo, err error) {
	err = r.Q.EventBooking.WithContext(ctx).UnderlyingDB().Raw(` SELECT 
        MIN(b.start_date) as start_date,
        MAX(b.end_date) as end_date,
		SUM(bp.total_price) as total_price,
		SUM(bp.paid) as total_paid,
		SUM(bp.discount) as discount
      FROM r_events as e
      JOIN r_booking_events as ref ON ref.event_id = e.id
      JOIN r_bookings as b ON b.id = ref.booking_id
	  JOIN r_booking_prices as bp on bp.id = b.id
      WHERE e.id = ?`, eventBookingID).Scan(&res).Error
	return
}

func (r *eventBookingRepository) CreateEventBooking(req *common.RequestContext, i dto.EventBookingData) (
	res dto.EventBookingDto, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var eventBooking model.EventBooking
	partyID, err := tx.EventBooking.WithContext(req.Ctx).InsertParty(proto.RegatePartyType_eventBooking.String())
	if err != nil {
		return res, err
	}
	fields := i.Fields
	eventBooking.ID = partyID
	eventBooking.CompanyID = req.ActiveCompany.ID
	eventBooking.Status = proto.State_ENABLED.String()
	if err = r.convertor.CopyStructData(fields, &eventBooking); err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).EventBooking.Save(&eventBooking)
	if err != nil {
		return
	}
	
	if err = tx.Deal.InsertActivity(partyID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	res = dto.EventBookingDtoFromModel(&eventBooking)
	err = tx.Commit()
	return
}
