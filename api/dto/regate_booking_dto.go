package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	RequestBookings struct {
		EventID    string `query:"event_id" required:"false"`
		CustomerID string `query:"customer_id" required:"false"`
		CourtID    string `query:"court_id" required:"false"`
		PaginationParams
		OptionalQueryParams
		AuthParams
	}

	UpdateBookingBatchRequest struct {
		Body struct {
			BookingIds  []int64 `json:"booking_ids"`
			TargetState string  `json:"target_state" minLength:"1"`
		}
	}

	CreateBookingRequest struct {
		Body CreateBookingBody
	}

	CreateBookingBody struct {
		CustomerID     int64         `json:"customer_id"`
		EventID        *int64         `json:"event_id" required:"false"`
		AdvancePayment float64       `json:"advance_payment" required:"false"`
		Comment        string        `json:"comment" required:"false"`
		Bookings       []BookingData `json:"bookings"`
	}

	RescheduleBooking struct {
		CourtID       int64     `json:"court_id" required:"true"`
		StartDateTime time.Time `json:"start_date" required:"true"`
		EndDateTime   time.Time `json:"end_date" required:"true"`
		IsValid       bool      `json:"is_valid"`

		Times      []string `json:"times"`
		DayWeek    int32    `json:"day_week"`
		TotalPrice int32    `json:"total_price" required:"false"`
	}

	BookingPaymentRequest struct {
		Body BookingPaymentBody
	}

	BookingPaymentBody struct {
		BookingID       int64   `json:"booking_id"`
		AddedAmount     float64 `json:"added_amount"`
		TotalPaidAmount float64 `json:"total_paid_amount"`
	}

	BookingRescheduleRequest struct {
		Body BookingRescheduleBody
	}

	BookingRescheduleBody struct {
		BookingID   int64       `json:"booking_id" required:"true"`
		BookingCode string      `json:"booking_code" required:"true"`
		PaidAmount  int32       `json:"paid_amount" required:"true"`
		PartyID     int64       `json:"party_id" required:"true"`
		BookingData BookingData `json:"booking" required:"true"`
	}

	ValidateBookingRequest struct {
		Body ValidateBookingData
	}

	ValidateBookingData struct {
		BookingID int64         `json:"booking_id" required:"false"`
		Bookings  []BookingData `json:"bookings" required:"true"`

		CustomerID   *int64  `json:"customer_id" required:"false"`
		CustomerName *string `json:"customer_name" required:"false"`
		EventID      *int64  `json:"event_id" required:"false"`
		EventName    *string `json:"event_name" required:"false"`
	}

	BookingData struct {
		AvailableCourts []AvailableCourtDto `json:"available_courts" required:"false"`
		CourtID         int64               `json:"court_id" required:"true"`
		CourtName       string              `json:"court_name" required:"true"`
		StartDateTime   time.Time           `json:"start_date" required:"true"`
		EndDateTime     time.Time           `json:"end_date" required:"true"`
		IsReserved      bool                `json:"is_reserved"`
		Times           []string            `json:"times"`
		DayWeek         int32               `json:"day_week"`
		TotalPrice      float64             `json:"total_price" required:"false"`
		Discount        float64             `json:"discount" required:"false"`
	}

	BookingDto struct {
		ID        int64     `json:"id"`
		Code      string    `json:"code"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		Type      string    `json:"-"`

		TotalPrice int64 `json:"total_price"`
		Paid       int64 `json:"paid"`
		Discount   int64 `json:"discount"`

		CourtName string `json:"court_name"`
		CourtUUID string `json:"court_uuid"`
		CourtID   int64  `json:"court_id"`

		PartyName string `json:"party_name"`
		PartyUUID string `json:"party_uuid"`
		PartyID   int64  `json:"party_id"`

		EventName *string `json:"evento_name"`
		EventID *int64 `json:"event_id"`

		Contacts []ContactDto `gorm:"-" json:"contacts"`
	}
)

func BookingDtoFromModel(m *model.Booking) BookingDto {
	return BookingDto{
		Code:      m.Code,
		Status:    m.Status,
		StartDate: m.StartDate,
		EndDate:   m.EndDate,
	}
}
