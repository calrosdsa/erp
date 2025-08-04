package booking_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"erp/pkg/di"
	activity_ucase "erp/project/core/activity/usecase"
	regate_event "erp/project/regate/internal/domain/event"
	"time"

	"gorm.io/gorm"
)

type BookingEventRepository interface {
	DeleteBookings(ctx context.Context, payload regate_event.BookingStatusData) (err error)
	CancelBookings(ctx context.Context, payload regate_event.BookingStatusData) (err error)
	CompletedBookings(ctx context.Context, payload regate_event.BookingStatusData) (err error)
	EditBookingPaidAmount(ctx context.Context, payload regate_event.EditPaidBookingEventData) (err error)
	OnRescheduleBooking(ctx context.Context, payload regate_event.RescheduleBookingEventData) (err error)
	OnEventCancelled(ctx context.Context, payload regate_event.StatusEventBookingData) (err error)
	OnEventCompleted(ctx context.Context, payload regate_event.StatusEventBookingData) (err error)
}

type bookingEventRepository struct {
	locale                helpers.Locale
	DB                    *gorm.DB
	convertor             helpers.ConvertorHelper
	currencyHelper        helpers.CurrencyHelper
	bookingSlorRepository BookingSlotRepository
	setting               repository.SettingService
	activityUcase         activity_ucase.ActivityUseCase
	c                     di.Container
}

func NewBookingEventRepo(
	helpers *helpers.Helpers,
	conn db.Connection,
	bookingSlorRepository BookingSlotRepository,
	setting repository.SettingService,
	c di.Container,
) BookingEventRepository {

	return &bookingEventRepository{
		locale:                helpers.Locale,
		convertor:             helpers.Convertor,
		currencyHelper:        helpers.Currency,
		DB:                    conn.GetDB(),
		bookingSlorRepository: bookingSlorRepository,
		setting:               setting,
		activityUcase:         c.Get(domain.ActivityUseCase).(activity_ucase.ActivityUseCase),
	}
}

func (r *bookingEventRepository) OnEventCompleted(ctx context.Context,
	payload regate_event.StatusEventBookingData) (err error) {
	event := payload.Event
	tx := payload.Tx
	bookingQ := tx.Booking
	bookingEventQ := tx.BookingEvent
	bookings, err := bookingQ.WithContext(ctx).
		Join(bookingEventQ, bookingEventQ.BookingID.EqCol(bookingQ.ID)).
		Where(
			bookingEventQ.EventID.Eq(event.ID),
			bookingQ.Status.Neq(proto.State_COMPLETED.String()),
			bookingQ.Status.Neq(proto.State_CANCELLED.String()),
		).
		Find()
	err = r.CompletedBookings(ctx, regate_event.BookingStatusData{
		Tx:              tx,
		Bookings:        bookings,
		CompanyDefaults: payload.CompanyDefaults,
		Profile:         payload.Profile,
	})
	return
}

func (r *bookingEventRepository) OnEventCancelled(ctx context.Context,
	payload regate_event.StatusEventBookingData) (err error) {
	event := payload.Event
	tx := payload.Tx
	bookingQ := tx.Booking
	bookingEventQ := tx.BookingEvent
	bookings, err := bookingQ.WithContext(ctx).Select(bookingQ.ID, bookingQ.Code).
		Join(bookingEventQ, bookingEventQ.BookingID.EqCol(bookingQ.ID)).
		Where(
			bookingEventQ.EventID.Eq(event.ID),
			bookingQ.Status.Neq(proto.State_CANCELLED.String()),
		).
		Find()
	err = r.CancelBookings(ctx, regate_event.BookingStatusData{
		Tx:              tx,
		Bookings:        bookings,
		CompanyDefaults: payload.CompanyDefaults,
		Profile:         payload.Profile,
	})
	return
}

func (r *bookingEventRepository) OnRescheduleBooking(ctx context.Context, payload regate_event.RescheduleBookingEventData,
) (err error) {
	body := payload.BookingReschedule
	tx := payload.Tx
	err = r.deleteBookingSlots(tx, ctx, []int64{body.BookingID})
	if err != nil {
		return
	}

	err = r.bookingSlorRepository.UpdatedCreatedBookings(tx, ctx, payload.Company.ID)
	if err != nil {
		return
	}
	return
}

func (r *bookingEventRepository) EditBookingPaidAmount(ctx context.Context,
	payload regate_event.EditPaidBookingEventData) (err error) {
	tx := payload.Tx
	booking := payload.Booking
	bookingPrice, err := tx.BookingPrice.WithContext(ctx).Where(
		tx.BookingPrice.ID.Eq(booking.ID),
	).First()
	if err != nil {
		return
	}
	//Update Booking slots
	err = r.updateBookingSlotPrices(tx, ctx, bookingPrice)
	if err != nil {
		return
	}
	//Update BookingState
	nextState := getBookingStateBaseOnAmount(bookingPrice.TotalPrice-bookingPrice.Discount, bookingPrice.Paid)
	_, err = tx.Booking.WithContext(ctx).Where(
		tx.Booking.ID.Eq(booking.ID),
	).UpdateSimple(tx.Booking.Status.Value(nextState))
	if err != nil {
		return
	}

	activity := model.Activity{
		ProfileID: payload.Profile.ID,
		PartyID:   payload.Booking.ID,
		Type:      proto.ActivityType_STAGE.String(),
	}
	r.activityUcase.CreateActivityStatus(tx, booking.Status, nextState, &activity)

	return
}

func getBookingStateBaseOnAmount(totalAmount, paidAmount int64) (bookingState string) {
	if paidAmount >= totalAmount {
		bookingState = proto.State_COMPLETED.String()
	} else if paidAmount > 0 {
		bookingState = proto.State_PARTIALLY_PAID.String()
	} else {
		bookingState = proto.State_UNPAID.String()
	}
	return
}
func (r *bookingEventRepository) DeleteBookings(ctx context.Context, payload regate_event.BookingStatusData) (err error) {
	tx := payload.Tx
	bookingQ := tx.Booking
	bookings := payload.Bookings
	bookingIds := make([]int64, len(bookings))
	for i, booking := range payload.Bookings {
		bookingIds[i] = booking.ID
	}
	_, err = bookingQ.WithContext(ctx).Where(
		bookingQ.ID.In(bookingIds...),
	).Delete()
	return
}

func (r *bookingEventRepository) CompletedBookings(ctx context.Context,
	payload regate_event.BookingStatusData) (
	err error) {
	tx := payload.Tx
	bookingQ := tx.Booking
	// status := proto.State_COMPLETED.String()
	bookingPriceQ := tx.BookingPrice
	bookingsIds := make([]int64, len(payload.Bookings))

	for i, booking := range payload.Bookings {
		bookingsIds[i] = booking.ID

		activity := &model.Activity{
			PartyID:   booking.ID,
			ProfileID: payload.Profile.ID,
			Type:      proto.ActivityType_STAGE.String(),
		}
		err = r.activityUcase.CreateActivityStatus(tx, booking.Status, proto.State_COMPLETED.String(), activity)
		if err != nil {
			return
		}
	}
	//Update state of bookings
	_, err = tx.Booking.WithContext(ctx).Where(
		bookingQ.ID.In(bookingsIds...),
	).UpdateSimple(bookingQ.Status.Value(proto.State_COMPLETED.String()))
	if err != nil {
		return
	}
	// //Update booking price
	_, err = tx.BookingPrice.WithContext(ctx).Where(
		bookingPriceQ.ID.In(bookingsIds...),
	).UpdateSimple(bookingPriceQ.Paid.SetCol(bookingPriceQ.TotalPrice.SubCol(bookingPriceQ.Discount)))
	if err != nil {
		return
	}

	err = r.updateBookingsSlotPrices(tx, ctx, payload.Bookings, payload.CompanyDefaults)
	if err != nil {
		return
	}

	return
}

func (r *bookingEventRepository) updateBookingsSlotPrices(tx *query.QueryTx,
	ctx context.Context, bookings []*model.Booking, companyDefaults model.CompanyDefault) (err error) {
	accountSettings, err := r.setting.GetAccountSettings(ctx, companyDefaults.CompanyID)
	if err != nil {
		return
	}
	var transactions []*model.TransactionLedger
	for _, booking := range bookings {
		bookingPrice, err := tx.BookingPrice.WithContext(ctx).Where(
			tx.BookingPrice.ID.Eq(booking.ID),
		).First()
		if err != nil {
			return err
		}
		debitTx := &model.TransactionLedger{
			Ledger:        accountSettings.CashAccunt,
			LedgerAgainst: &accountSettings.ReceivableAccount,
			Debit:         int64(bookingPrice.Paid),
			Currency:      companyDefaults.Currency,
			VoucherCode:   booking.Code,
			VoucherType:   proto.RegatePartyType_booking.String(),
			PostingDate:   time.Now(),
		}
		creditTx := &model.TransactionLedger{
			Ledger:        accountSettings.ReceivableAccount,
			LedgerAgainst: &accountSettings.CashAccunt,
			Credit:        int64(bookingPrice.Paid),
			Currency:      companyDefaults.Currency,
			VoucherCode:   booking.Code,
			VoucherType:   proto.RegatePartyType_booking.String(),
			PartyID:       &booking.Party,
			PostingDate:   time.Now(),
		}
		transactions = append(transactions, debitTx, creditTx)
		err = r.updateBookingSlotPrices(tx, ctx, bookingPrice)
		if err != nil {
			return err
		}
	}
	err = tx.TransactionLedger.WithContext(ctx).CreateInBatches(transactions, len(transactions))
	return
}

func (r *bookingEventRepository) updateBookingSlotPrices(tx *query.QueryTx, ctx context.Context,
	bookingPrice *model.BookingPrice) (err error) {
	totalPaidAmount := bookingPrice.Paid
	bookingSlotQ := tx.BookingSlot
	bookingSlots, err := bookingSlotQ.WithContext(ctx).
		Select(bookingSlotQ.ID, bookingSlotQ.PaidAmount, bookingSlotQ.TotalPrice).
		Where(bookingSlotQ.BookingID.Eq(bookingPrice.ID)).Order(
		bookingSlotQ.Datetime.Asc(),
	).Find()
	for _, bookingSlot := range bookingSlots {
		bookingSlotPaidAmount := bookingSlot.PaidAmount
		if totalPaidAmount > bookingSlot.TotalPrice {
			bookingSlotPaidAmount = bookingSlot.TotalPrice
			totalPaidAmount -= bookingSlot.TotalPrice
		} else if totalPaidAmount > 0 {
			bookingSlotPaidAmount = totalPaidAmount
			totalPaidAmount = 0

		} else {
			bookingSlotPaidAmount = 0
		}
		if bookingSlotPaidAmount != bookingSlot.PaidAmount {
			bookingSlotQ.Where(
				bookingSlotQ.ID.Eq(bookingSlot.ID),
			).UpdateSimple(bookingSlotQ.PaidAmount.Value(bookingSlotPaidAmount))
		}

	}

	return
}

func (r *bookingEventRepository) CancelBookings(ctx context.Context, payload regate_event.BookingStatusData) (
	err error) {
	tx := payload.Tx
	bookingQ := tx.Booking
	// targetStatus := proto.State_CANCELLED.String()
	bookingsIds := make([]int64, len(payload.Bookings))

	for i, booking := range payload.Bookings {
		bookingsIds[i] = booking.ID
		activity := &model.Activity{
			PartyID:   booking.ID,
			ProfileID: payload.Profile.ID,
			Type:      proto.ActivityType_STAGE.String(),
		}
		err = r.activityUcase.CreateActivityStatus(tx, booking.Status, proto.State_CANCELLED.String(), activity)
		if err != nil {
			return
		}
		//Delete Transaction ledgers
		_, err = tx.TransactionLedger.WithContext(ctx).Where(
			tx.TransactionLedger.VoucherCode.Eq(booking.Code),
		).Delete()
		if err != nil {
			return
		}
	}
	//Update state of bookings
	_, err = tx.Booking.WithContext(ctx).Where(
		bookingQ.ID.In(bookingsIds...),
	).UpdateSimple(bookingQ.Status.Value(proto.State_CANCELLED.String()))
	if err != nil {
		return
	}
	//Delete booking slots
	err = r.deleteBookingSlots(tx, ctx, bookingsIds)
	if err != nil {
		return
	}
	return
}

func (r *bookingEventRepository) deleteBookingSlots(tx *query.QueryTx, ctx context.Context, bookingsIds []int64) error {
	_, err := tx.BookingSlot.WithContext(ctx).Where(
		tx.BookingSlot.BookingID.In(bookingsIds...),
	).Delete()
	return err
}
