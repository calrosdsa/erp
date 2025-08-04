package dto

import "time"

type (
	RequestBookingSlots struct {
		CourtID  string    `query:"court_id" required:"true"`
		FromDate time.Time `query:"from_date" required:"true"`
		ToDate   time.Time `query:"to_date" required:"true"`
	}
	BookingSlotDto struct {
		TotalPrice int32     `json:"total_price"`
		PaidPrice  int32     `json:"paid_price"`
		BookingID  int64     `json:"booking_id"`
		Datetime   time.Time `json:"datetime"`
		Type       string    `json:"type"`

		PartyName   string `json:"party_name"`
		BookingCode string `json:"booking_code"`
	}

	BookingScheduleBody struct {
		Body struct {
			BookingSlots []BookingSlotDto `json:"booking_slots"`
			CourtRates   []CourtRateDto   `json:"cour_rates"`
		}
	}
)
