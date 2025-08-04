package booking_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type BookingRepository interface {
	BookingReschedule(tx *query.QueryTx, req *common.RequestContext, d dto.BookingRescheduleBody) error
	CreateBooking(tx *query.QueryTx, req *common.RequestContext, d dto.CreateBookingBody) error
	GetBookings(req *common.RequestContext, i *dto.RequestBookings) (
		res dto.PaginationResult[[]dto.BookingDto], err error)
	GetBooking(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.BookingDto], err error)
	ValidateBooking(req *common.RequestContext, i dto.ValidateBookingData) (
		res dto.ValidateBookingData, err error)
	UpdateBookingState(tx *query.QueryTx, req *common.RequestContext, id string, prevState, nexState string) (
		*model.Booking, error)
	EditBookingPaidAmount(tx *query.QueryTx, req *common.RequestContext, d dto.BookingPaymentBody) (
		*model.Booking, error)
	UpdateBookingBatch(req *common.RequestContext, d *dto.UpdateBookingBatchRequest) ([]*model.Booking, error)
}

type bookingRepository struct {
	conn           db.Connection
	Q              *query.Query
	DB             *gorm.DB
	currencyHelper helpers.CurrencyHelper
	convertor      helpers.ConvertorHelper
	setting        repository.SettingService
	accounting     repository.AccountingService
}

func NewBookingRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
	settingS repository.SettingService,
	accounting repository.AccountingService,
) BookingRepository {
	return &bookingRepository{
		conn:           conn,
		Q:              conn.GetQ(),
		DB:             conn.GetDB(),
		currencyHelper: helpers.Currency,
		convertor:      helpers.Convertor,
		setting:        settingS,
		accounting:     accounting,
	}
}

func (r *bookingRepository) UpdateBookingBatch(req *common.RequestContext, d *dto.UpdateBookingBatchRequest) (
	res []*model.Booking, err error) {
	res, err = r.Q.Booking.WithContext(req.Ctx).Where(
		r.Q.Booking.ID.In(d.Body.BookingIds...),
	).Find()
	return
}

func (r *bookingRepository) BookingReschedule(tx *query.QueryTx, req *common.RequestContext,
	d dto.BookingRescheduleBody) (err error) {
	totalPrice :=r.currencyHelper.FloatToInt(d.BookingData.TotalPrice)
	discount :=r.currencyHelper.FloatToInt(d.BookingData.Discount)
	paidAmount := d.PaidAmount
	bookingQ := tx.Booking
	bookingPriceQ := tx.BookingPrice
	_, err = bookingQ.WithContext(req.Ctx).Where(
		bookingQ.CompanyID.Eq(req.ActiveCompany.ID),
		bookingQ.ID.Eq(d.BookingID),
	).UpdateSimple(
		bookingQ.StartDate.Value(d.BookingData.StartDateTime),
		bookingQ.EndDate.Value(d.BookingData.EndDateTime),
		bookingQ.CourtID.Value(d.BookingData.CourtID),
		bookingQ.Status.Value(proto.State_DRAFT.String()),
	)
	if err != nil {
		return
	}

	_, err = bookingPriceQ.WithContext(req.Ctx).Where(
		bookingPriceQ.ID.Eq(d.BookingID),
	).UpdateSimple(
		bookingPriceQ.TotalPrice.Value(totalPrice),
		bookingPriceQ.Discount.Value(discount),
	)
	err = r.accounting.DelTxnsByVoucherCode(req.Ctx, tx, d.BookingCode)
	if err != nil {
		return
	}
	accountSettings, err := r.setting.GetAccountSettings(req.Ctx, req.ActiveCompany.ID)
	if err != nil {
		return
	}
	transactions := r.bookingTxLedger(int64(totalPrice), int64(discount), int64(paidAmount), req.CompanyDefaults.Currency,
		d.PartyID, accountSettings.CashAccunt, accountSettings.ReceivableAccount, accountSettings.IncomeAccount, d.BookingCode)
	err = tx.TransactionLedger.WithContext(req.Ctx).CreateInBatches(transactions, len(transactions))
	if err != nil {
		return
	}
	return
}

func (r *bookingRepository) EditBookingPaidAmount(tx *query.QueryTx, req *common.RequestContext, d dto.BookingPaymentBody) (
	res *model.Booking, err error) {
	bookingQ := tx.Booking
	bookingPriceQ := tx.BookingPrice
	addedAmount := r.currencyHelper.FloatToInt(d.AddedAmount)
	totalPaidAmount := r.currencyHelper.FloatToInt(d.TotalPaidAmount)
	res, err = bookingQ.WithContext(req.Ctx).Where(
		bookingQ.ID.Eq(d.BookingID),
		bookingQ.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}
	_, err = bookingPriceQ.WithContext(req.Ctx).Where(
		bookingPriceQ.ID.Eq(res.ID),
	).UpdateSimple(bookingPriceQ.Paid.Value(
		totalPaidAmount,
	))
	if err != nil {
		return
	}
	accountSettings, err := r.setting.GetAccountSettings(req.Ctx, req.ActiveCompany.ID)
	if err != nil {
		return
	}
	activityData := map[string]interface{}{
		"title":       "Se ha editado el monto pagado",
		"description": fmt.Sprintf(`Monto pagado de "%.2f" a "%.2f"`, d.TotalPaidAmount-d.AddedAmount, d.TotalPaidAmount),
	}
	jsonString := domain.JsonStringify(activityData)
	activity := &model.Activity{
		ProfileID: req.Profile.ID,
		Type:      proto.ActivityType_INFO.String(),
		PartyID:   d.BookingID,
		Data:      &jsonString,
	}
	err = tx.WithContext(req.Ctx).Activity.Save(activity)
	if err != nil {
		return
	}

	transactions := make([]*model.TransactionLedger, 2)
	debitTx := &model.TransactionLedger{
		Ledger:        accountSettings.CashAccunt,
		LedgerAgainst: &accountSettings.ReceivableAccount,
		Debit:         addedAmount,
		Currency:      req.CompanyDefaults.Currency,
		VoucherCode:   res.Code,
		VoucherType:   proto.RegatePartyType_booking.String(),
		PostingDate:   time.Now(),
	}
	creditTx := &model.TransactionLedger{
		Ledger:        accountSettings.ReceivableAccount,
		LedgerAgainst: &accountSettings.CashAccunt,
		Credit:        addedAmount,
		Currency:      req.CompanyDefaults.Currency,
		VoucherCode:   res.Code,
		VoucherType:   proto.RegatePartyType_booking.String(),
		PartyID:       &res.Party,
		PostingDate:   time.Now(),
	}
	transactions[0] = debitTx
	transactions[1] = creditTx
	err = tx.WithContext(req.Ctx).TransactionLedger.CreateInBatches(transactions, domain.DEFAULT_BATCH_SIZE)
	return
}

func (r *bookingRepository) UpdateBookingState(tx *query.QueryTx, req *common.RequestContext, id string,
	prevState, nexState string) (res *model.Booking, err error) {
	bookingID := r.convertor.StrtoInt(id)
	bookingQ := tx.Booking
	// _, err = tx.Booking.WithContext(req.Ctx).Where(
	// 	bookingQ.ID.Eq(bookingID),
	// 	bookingQ.CompanyID.Eq(req.ActiveCompany.ID)
	// 	bookingQ.Status.Eq(prevState),
	// ).UpdateSimple(bookingQ.Status.Value(nexState))
	// if err != nil {
	// 	return
	// }
	booking, err := tx.Booking.WithContext(req.Ctx).Where(
		bookingQ.ID.Eq(bookingID),
	).First()
	if err != nil {
		return
	}
	return booking, err
}

func (r *bookingRepository) CreateBooking(tx *query.QueryTx, req *common.RequestContext, d dto.CreateBookingBody) (
	err error) {
	bookingType := proto.BookingType_LOCAL_BOOKING.String()
	if d.EventID != nil{
		bookingType = proto.BookingType_EVENT_BOOKING.String()
	}
	bookings := make([]*model.Booking, len(d.Bookings))
	bookingsActivities := make([]*model.Activity, len(d.Bookings))
	bookingPrices := make([]*model.BookingPrice, len(d.Bookings))

	var transactions []*model.TransactionLedger

	accountSettings, err := r.setting.GetAccountSettings(req.Ctx, req.ActiveCompany.ID)
	if err != nil {
		return
	}
	activityComments := make([]*model.ActivityComment, len(d.Bookings))

	bookingEventReferences := make([]*model.BookingEvent, len(d.Bookings))

	bookingPartyReferences := make([]*model.PartyReference, len(d.Bookings))

	advancePayment := r.currencyHelper.FloatToInt(d.AdvancePayment)

	for i, booking := range d.Bookings {
		m := &model.Booking{}
		bookingID, err1 := tx.Booking.InsertParty(proto.RegatePartyType_booking.String())
		if err1 != nil {
			err = err1
			return
		}
		m.Code = r.conn.GenerateCode(req.Ctx, model.Booking{}, req.ActiveCompany.ID)
		m.ID = bookingID
		m.StartDate = booking.StartDateTime
		m.EndDate = booking.EndDateTime
		m.CourtID = booking.CourtID
		m.Party = d.CustomerID
		m.CompanyID = req.ActiveCompany.ID
		m.Status = proto.State_DRAFT.String()
		m.Type = bookingType
		bookings[i] = m

		bookingPrice := &model.BookingPrice{}
		bookingPrice.ID = m.ID
		totalPrice := r.currencyHelper.FloatToInt(booking.TotalPrice)
		discount := r.currencyHelper.FloatToInt(booking.Discount)
		bookingPrice.TotalPrice = totalPrice
		bookingPrice.Discount = discount
		if advancePayment > 0 {
			if advancePayment > int64(totalPrice) {
				bookingPrice.Paid = totalPrice
			} else {
				bookingPrice.Paid = advancePayment
			}
			advancePayment -= int64(totalPrice)
		}
		bookingPrices[i] = bookingPrice
		txLedger := r.bookingTxLedger(int64(bookingPrice.TotalPrice), int64(bookingPrice.Discount), int64(bookingPrice.Paid),
			req.CompanyDefaults.Currency, m.Party, accountSettings.CashAccunt, accountSettings.ReceivableAccount,
			accountSettings.IncomeAccount, m.Code)

		transactions = append(transactions, txLedger...)
		//Add Booking event reference
		if d.EventID != nil {
			eventReference := &model.BookingEvent{}
			eventReference.BookingID = bookingID
			eventReference.EventID = *d.EventID
			bookingEventReferences[i] = eventReference
		}
		//Add party reference
		partyReference := &model.PartyReference{}
		partyReference.PartyID = bookingID
		partyReference.ReferenceID = d.CustomerID
		bookingPartyReferences[i] = partyReference
		//add booking comment
		if d.Comment != "" {
			activityComment := &model.Activity{}
			activityComment.PartyID = bookingID
			activityComment.ProfileID = req.Profile.ID
			activityComment.Type = proto.ActivityType_COMMENT.String()
			err = tx.Activity.Save(activityComment)
			if err != nil {
				return
			}
			// activityComment.Comment = &d.Comment

			comment := &model.ActivityComment{}
			comment.Comment = d.Comment
			comment.ActivityID = activityComment.ID
			activityComments[i] = comment

		}

		//Add Activity booking
		activity := &model.Activity{}
		activity.Type = proto.ActivityType_CREATE.String()
		activity.PartyID = bookingID
		activity.ProfileID = req.Profile.ID
		bookingsActivities[i] = activity
	}

	//creating bookings
	err = tx.Booking.WithContext(req.Ctx).CreateInBatches(bookings, len(bookings))
	if err != nil {
		return
	}
	//creating bookings price
	err = tx.BookingPrice.WithContext(req.Ctx).CreateInBatches(bookingPrices, len(bookingPrices))
	if err != nil {
		return
	}
	//create bookings events
	if d.EventID != nil {
		err = tx.BookingEvent.WithContext(req.Ctx).CreateInBatches(bookingEventReferences,
			len(bookingEventReferences))
		if err != nil {
			return
		}
	}
	//create customer references
	err = tx.PartyReference.WithContext(req.Ctx).CreateInBatches(bookingPartyReferences, len(bookingPartyReferences))
	if err != nil {
		return
	}
	//creaate transactin ledgers
	err = tx.TransactionLedger.WithContext(req.Ctx).CreateInBatches(transactions, len(transactions))
	if err != nil {
		return
	}
	if d.Comment != "" {
		err = tx.ActivityComment.WithContext(req.Ctx).CreateInBatches(activityComments, len(activityComments))
		if err != nil {
			return
		}
	}

	//creating bookings activity
	err = tx.Activity.WithContext(req.Ctx).CreateInBatches(bookingsActivities, len(bookingsActivities))

	return
}

func (r *bookingRepository) GetBookings(req *common.RequestContext, i *dto.RequestBookings) (
	res dto.PaginationResult[[]dto.BookingDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	booking := r.Q.Booking
	builder := r.Q.WithContext(req.Ctx).Booking

	//ADDING CONDITIONS

	customerQ := r.Q.Customer
	courtQ := r.Q.Court
	eventReferenceQ := r.Q.BookingEvent
	bookingPriceQ := r.Q.BookingPrice
	builder = builder.Select(
		booking.ID, booking.Code, booking.Status, booking.Type, booking.CreatedAt,
		booking.StartDate, booking.EndDate,
		booking.Party.As("party_id"), customerQ.Name.As("party_name"), customerQ.UUID.As("party_uuid"),
		booking.CourtID, courtQ.Name.As("court_name"), courtQ.UUID.As("court_uuid"),
		bookingPriceQ.TotalPrice, bookingPriceQ.Paid, bookingPriceQ.Discount,
	).
		Join(customerQ, customerQ.ID.EqCol(booking.Party)).
		Join(courtQ, courtQ.ID.EqCol(booking.CourtID)).
		Join(bookingPriceQ, bookingPriceQ.ID.EqCol(booking.ID))

	conds = append(conds, booking.CompanyID.Eq(req.ActiveCompany.ID))

	if i.EventID != "" {
		eventID := r.convertor.StrtoInt(i.EventID)
		builder = builder.Join(eventReferenceQ, eventReferenceQ.BookingID.EqCol(booking.ID))
		conds = append(conds, eventReferenceQ.EventID.Eq(eventID))
	}
	if i.CourtID != "" {
		courtID := r.convertor.StrtoInt(i.CourtID)
		conds = append(conds, booking.CourtID.Eq(courtID))
	}
	if i.CustomerID != "" {
		customerID := r.convertor.StrtoInt(i.CustomerID)
		conds = append(conds, booking.Party.Eq(customerID))
	}
	if i.Status != "" {
		conds = append(conds, booking.Status.Eq(i.Status))
	}

	//Filtering
	builder = builder.Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}

	//Pagination Order
	limit, offset := r.convertor.ToPaginationParams(i.Page, i.Size)
	orderCol, ok := r.Q.Booking.GetFieldByName(i.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if strings.ToUpper(i.Order) == "ASC" {
			fmt.Println("ORDER", "ASC")
			order = orderCol.Asc()
		} else {
			fmt.Println("ORDER", "DESC")
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}

func (r *bookingRepository) GetBooking(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.BookingDto], err error) {
	bookingQ := r.Q.Booking
	customerQ := r.Q.Customer
	courtQ := r.Q.Court
	bookingPriceQ := r.Q.BookingPrice
	eventBookingQ := r.Q.EventBooking
	bookingEventQ := r.Q.BookingEvent
	bookingID := r.convertor.StrtoInt(i.ID)
	err = r.Q.Booking.WithContext(req.Ctx).Select(
		bookingQ.ID, bookingQ.Code, bookingQ.Status, bookingQ.Type, bookingQ.CreatedAt,
		bookingQ.StartDate, bookingQ.EndDate, bookingQ.Party.As("party_id"),
		customerQ.Name.As("party_name"), customerQ.UUID.As("party_uuid"),
		courtQ.Name.As("court_name"), courtQ.UUID.As("court_uuid"), courtQ.ID.As("court_id"),
		bookingPriceQ.TotalPrice, bookingPriceQ.Paid, bookingPriceQ.Discount,
		eventBookingQ.Name.As("event_name"), eventBookingQ.ID.As("event_id"),
	).
		Join(customerQ, customerQ.ID.EqCol(bookingQ.Party)).
		Join(courtQ, courtQ.ID.EqCol(bookingQ.CourtID)).
		Join(bookingPriceQ, bookingPriceQ.ID.EqCol(bookingQ.ID)).
		LeftJoin(bookingEventQ, bookingQ.Type.Eq(proto.BookingType_EVENT_BOOKING.String()),
			bookingEventQ.BookingID.EqCol(bookingQ.ID)).
		LeftJoin(eventBookingQ, eventBookingQ.ID.EqCol(bookingEventQ.EventID)).
		Where(bookingQ.CompanyID.Eq(req.ActiveCompany.ID), bookingQ.ID.Eq(bookingID)).Scan(&res.Entity)
	if err != nil {
		return
	}

	return
}
func (r *bookingRepository) ValidateBooking(req *common.RequestContext, d dto.ValidateBookingData) (
	dto.ValidateBookingData, error) {
	bookingSlotQ := r.Q.BookingSlot
	courtRateQ := r.Q.CourtRate
	res := make([]dto.BookingData, len(d.Bookings))
	for i, booking := range d.Bookings {
		fmt.Println("DATETIME", booking.StartDateTime, booking.EndDateTime, d.BookingID)
		bookingData := booking
		// Checking for existing conflicting bookings
		builder := r.Q.BookingSlot.WithContext(req.Ctx)
		if d.BookingID != 0 {
			fmt.Println("ADDING NOT CONDITION")
			builder = builder.Not(bookingSlotQ.BookingID.Eq(d.BookingID))

		}
		builder = builder.Where(
			bookingSlotQ.CourtID.Eq(booking.CourtID),
			bookingSlotQ.Datetime.Gt(booking.StartDateTime),
			bookingSlotQ.Datetime.Lt(booking.EndDateTime),
		)

		// Count the results
		count, err := builder.Count()
		if err != nil {
			return dto.ValidateBookingData{}, err
		}
		fmt.Println("count bookings", count)
		if count > 0 {
			bookingData.IsReserved = true
			bookingData.AvailableCourts, err = r.getAvailableCourts(req, booking)
			if err != nil {
				fmt.Println("AVAILABLE COURTS ERROR", err)
			}
		}
		fmt.Println(bookingData.IsReserved)
		fmt.Println("TIMES", booking.Times)

		courtRates, err := courtRateQ.WithContext(req.Ctx).Where(
			courtRateQ.CourtID.Eq(booking.CourtID),
			courtRateQ.Enabled.Is(true),
			courtRateQ.DayWeek.Eq(booking.DayWeek),
			courtRateQ.Time.In(booking.Times...),
		).Select(courtRateQ.Rate).Find()
		if err != nil {
			return dto.ValidateBookingData{}, err
		}
		sumTotal := r.sumRate(courtRates)
		bookingData.TotalPrice = r.currencyHelper.IntToFloat(int64(sumTotal))
		fmt.Println("TOTAL PRICE", bookingData.TotalPrice)
		fmt.Println("COURT RATES", courtRates)
		res[i] = bookingData
		d.Bookings = res
	}
	return d, nil
}

func (r *bookingRepository) getAvailableCourts(req *common.RequestContext, bookingData dto.BookingData) (
	res []dto.AvailableCourtDto, err error,
) {
	query := `SELECT c.*,
       COALESCE(
           (SELECT SUM(rate)
            FROM r_court_rates
            WHERE time = ANY($4)
              AND enabled = TRUE
              AND court_id = c.id
              AND day_week = $5), 0) AS total_price
			FROM r_courts AS c
			WHERE company_id = $1 
			AND enabled = TRUE
			AND NOT EXISTS (
				SELECT 1
				FROM r_bookings
				WHERE court_id = c.id
					AND (
						($2 BETWEEN start_date AND end_date)
						OR 
						($3 BETWEEN start_date AND end_date)
					)
			);
		`

	err = r.DB.WithContext(req.Ctx).Raw(
		query, req.ActiveCompany.ID, bookingData.StartDateTime, bookingData.EndDateTime,
		pq.Array(bookingData.Times), bookingData.DayWeek,
	).Limit(domain.DEFAULT_LIMIT).Scan(&res).Error

	return
}

func (r *bookingRepository) bookingTxLedger(totalAmount, discountAmount, paidAmount int64, currency string,
	partyID, cashAccount, receivableAccount, incomeAccount int64, voucherCode string) []*model.TransactionLedger {
	var transactions []*model.TransactionLedger
	fmt.Println("CASH ACCOUNT", cashAccount)
	fmt.Println("RECEIVABLE ACCOUNT", receivableAccount)
	fmt.Println("INCOME ACCCOUNT", incomeAccount)

	if paidAmount >= totalAmount {
		debitTx := &model.TransactionLedger{
			Ledger: cashAccount,
			// LedgerAgainst: &accountSettings.IncomeAccount,
			Debit:       paidAmount,
			Currency:    currency,
			VoucherCode: voucherCode,
			VoucherType: proto.RegatePartyType_booking.String(),
			PostingDate: time.Now(),
		}
		creditTx := &model.TransactionLedger{
			Ledger: incomeAccount,
			// LedgerAgainst: &cashAccount,
			Credit:      paidAmount,
			Currency:    currency,
			VoucherCode: voucherCode,
			VoucherType: proto.RegatePartyType_booking.String(),
			PartyID:     &partyID,
			PostingDate: time.Now(),
		}
		transactions = append(transactions, debitTx, creditTx)
	} else {
		debitTx := &model.TransactionLedger{
			Ledger: receivableAccount,
			// LedgerAgainst: &incomeAccount,
			Debit:       totalAmount - discountAmount,
			Currency:    currency,
			VoucherCode: voucherCode,
			VoucherType: proto.RegatePartyType_booking.String(),
			PartyID:     &partyID,
			PostingDate: time.Now(),
		}
		creditTx := &model.TransactionLedger{
			Ledger: incomeAccount,
			// LedgerAgainst: &receivableAccount,
			Credit:      totalAmount - discountAmount,
			Currency:    currency,
			VoucherCode: voucherCode,
			VoucherType: proto.RegatePartyType_booking.String(),
			PostingDate: time.Now(),
		}
		transactions = append(transactions, debitTx, creditTx)
		if paidAmount > 0 {
			debitTx := &model.TransactionLedger{
				Ledger: cashAccount,
				// LedgerAgainst: &receivableAccount,
				Debit:       paidAmount,
				Currency:    currency,
				VoucherCode: voucherCode,
				VoucherType: proto.RegatePartyType_booking.String(),
				PostingDate: time.Now(),
			}
			creditTx := &model.TransactionLedger{
				Ledger: receivableAccount,
				// LedgerAgainst: &cashAccount,
				Credit:      paidAmount,
				Currency:    currency,
				VoucherCode: voucherCode,
				VoucherType: proto.RegatePartyType_booking.String(),
				PartyID:     &partyID,
				PostingDate: time.Now(),
			}
			transactions = append(transactions, debitTx, creditTx)
		}
	}
	return transactions
}

func (r *bookingRepository) sumRate(d []*model.CourtRate) (res int64) {
	for _, courtRate := range d {
		fmt.Println("COURT RATE", courtRate.Rate)
		res += courtRate.Rate
	}
	return
}
