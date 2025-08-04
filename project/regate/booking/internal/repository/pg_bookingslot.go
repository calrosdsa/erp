package booking_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	regate_domain "erp/project/regate/internal/domain"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BookingSlotRepository interface {
	UpdatedCreatedBookings(tx *query.QueryTx, ctx context.Context, companyID int64) error
	GetBookingSlots(req *common.RequestContext, d *dto.RequestBookingSlots) (
		res dto.BookingScheduleBody, err error)
}

type bookingSlotRepository struct {
	conn      db.Connection
	Q         *query.Query
	DB        *gorm.DB
	convertor helpers.ConvertorHelper
}

func NewBookingSlotRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) BookingSlotRepository {
	return &bookingSlotRepository{
		conn:      conn,
		Q:         conn.GetQ(),
		DB:        conn.GetDB(),
		convertor: helpers.Convertor,
	}
}
func (r *bookingSlotRepository) GetBookingSlots(req *common.RequestContext, d *dto.RequestBookingSlots) (
	res dto.BookingScheduleBody, err error) {
	bookingSlotQ := r.Q.BookingSlot
	customerQ := r.Q.Customer
	bookingQ := r.Q.Booking
	courtID := r.convertor.StrtoInt(d.CourtID)
	fmt.Println("DATE",d.FromDate,d.ToDate)
	err = bookingSlotQ.WithContext(req.Ctx).Where(
		bookingSlotQ.Datetime.Between(d.FromDate, d.ToDate),
		bookingSlotQ.CourtID.Eq(courtID),
	).Select(
		bookingSlotQ.TotalPrice, bookingSlotQ.PaidAmount, bookingSlotQ.BookingID,
		bookingSlotQ.Datetime, bookingSlotQ.Type,
		customerQ.Name.As("party_name"),
		bookingQ.Code.As("booking_code"),
	).
		Join(bookingQ, bookingQ.ID.EqCol(bookingSlotQ.BookingID)).
		Join(customerQ, customerQ.ID.EqCol(bookingQ.Party)).
		Scan(&res.Body.BookingSlots)
	if err != nil {
		return
	}
	fmt.Println("LEN",len(res.Body.BookingSlots))
	courRates, err := r.getCourtRates(req, courtID)
	if err != nil {
		return
	}
	res.Body.CourtRates = courRates
	return
}
func (r *bookingSlotRepository) getCourtRates(req *common.RequestContext, courtID int64) (
	res []dto.CourtRateDto, err error) {
	courtRateQ := r.Q.CourtRate
	err = r.Q.CourtRate.WithContext(req.Ctx).Select(
		courtRateQ.Rate, courtRateQ.Enabled, courtRateQ.Currency, courtRateQ.DayWeek, courtRateQ.Time,
		courtRateQ.Currency,
	).Where(courtRateQ.CourtID.Eq(courtID)).Scan(&res)
	if err != nil {
		return
	}

	return
}

func (r *bookingSlotRepository) UpdatedCreatedBookings(tx *query.QueryTx, ctx context.Context, companyID int64) (err error) {
	bookingQ := tx.Booking
	bookingPrice := tx.BookingPrice
	var res []dto.BookingDto
	err = tx.Booking.WithContext(ctx).Select(
		bookingQ.CourtID, bookingQ.Status, bookingQ.ID, bookingQ.StartDate, bookingQ.EndDate, bookingQ.Type,
		bookingPrice.TotalPrice, bookingPrice.Paid, bookingPrice.Discount,
	).
		Join(bookingPrice, bookingPrice.ID.EqCol(bookingQ.ID)).
		Where(
			bookingQ.CompanyID.Eq(companyID),
			bookingQ.Status.Eq(proto.State_DRAFT.String()),
		).Scan(&res)
	if err != nil {
		return
	}
	var slots []*model.BookingSlot
	for _, bookingData := range res {

		bookingSlots, err1 := r.insertBookingSlots(ctx, bookingData, companyID)
		if err1 != nil {
			return err1
		}
		fmt.Println("BOOKING SLOTS", bookingSlots)
		slots = append(slots, bookingSlots...)
		nextState := getBookingStateBaseOnAmount(bookingData.TotalPrice-bookingData.Discount, bookingData.Paid)
		_, err = tx.Booking.WithContext(ctx).Where(bookingQ.ID.Eq(bookingData.ID)).Update(bookingQ.Status, nextState)
		if err != nil {
			return
		}
	}
	err = tx.BookingSlot.WithContext(ctx).CreateInBatches(slots, len(slots))
	if err != nil {
		return
	}
	return
}

func (r *bookingSlotRepository) insertBookingSlots(ctx context.Context, d dto.BookingDto, companyID int64) (
	[]*model.BookingSlot, error) {
	var paidAmount int64
	bookingTimes := r.generateTimeSlots(d.StartDate, d.EndDate)
	if len(bookingTimes) == 0 {
		return []*model.BookingSlot{}, nil
	}
	startTime := bookingTimes[0]
	endTime := bookingTimes[len(bookingTimes)-1]
	fmt.Println("BOOKING TINE", bookingTimes, "BOOKING TYPE", d.Type)
	bookingSlots := make([]*model.BookingSlot, len(bookingTimes))
	//get court rates
	var courtRates []model.CourtRate
	query := `select coalesce(rate,0) AS rate from 
	generate_series($1 AT TIME ZONE $4,
	$2 AT TIME ZONE $4,interval '30 min') t(x)
    left join r_court_rates on court_id = $3 and enabled is true 
    and time = x::time 
	AND day_week = EXTRACT(dow FROM $1::date);`
	err := r.DB.WithContext(ctx).Raw(query, startTime, endTime, d.CourtID, regate_domain.DEFAULT_TIMEZONE).Scan(&courtRates).Error
	if err != nil {
		return bookingSlots, err
	}

	paidAmount = d.Paid
	var discountBySlot int64
	if len(bookingSlots) > 0 {
		discountBySlot = d.Discount / int64(len(bookingTimes))
	}

	for i, bookingTime := range bookingTimes {
		courtRate := courtRates[i].Rate
		fmt.Println("COURT RATE", courtRate)
		m := &model.BookingSlot{}
		m.BookingID = d.ID
		m.CompanyID = companyID
		m.CourtID = d.CourtID
		m.Datetime = bookingTime
		m.Type = d.Type
		fmt.Println("BOOKING SLOT TOTAL PRICE", m.TotalPrice)
		m.TotalPrice = courtRate - discountBySlot
		fmt.Println("Booking slot after discount", m.TotalPrice)
		if courtRate < paidAmount {
			m.PaidAmount = courtRate
			paidAmount -= courtRate
		} else if paidAmount > 0 {
			m.PaidAmount = paidAmount
			paidAmount = 0
		}
		bookingSlots[i] = m
	}
	return bookingSlots, nil
}

func (r *bookingSlotRepository) generateTimeSlots(startTime, endTime time.Time) []time.Time {
	// Slice to hold the time slots
	var timeSlots []time.Time

	// Generate time slots in 30-minute intervals
	for t := startTime; t.Before(endTime); t = t.Add(30 * time.Minute) {
		timeSlots = append(timeSlots, t)
	}

	return timeSlots
}
