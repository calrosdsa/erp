package regate_event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type BookingStatusData struct {
	Bookings        []*model.Booking
	Tx              *query.QueryTx
	Profile         model.Profile
	CompanyDefaults model.CompanyDefault
}

type OnCancelBookingEventData struct {
	Bookings []*model.Booking
	Profile  model.Profile
	Tx       *query.QueryTx
}

type EditPaidBookingEventData struct {
	Profile model.Profile
	Booking *model.Booking
	Tx      *query.QueryTx
}

type StatusEventBookingData struct {
	Event *model.EventBooking
	Tx *query.QueryTx
	Profile         model.Profile
	CompanyDefaults model.CompanyDefault
}

type RescheduleBookingEventData struct {
	Tx                *query.QueryTx
	BookingReschedule dto.BookingRescheduleBody
	Company           model.Company
}
