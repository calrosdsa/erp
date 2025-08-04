package internal

import (
	"gorm.io/gen"
)

func RegateModels(g *gen.Generator) []interface{} {
	court := g.GenerateModelAs("r_courts","Court")
	courtRate := g.GenerateModelAs("r_court_rates","CourtRate")
	booking := g.GenerateModelAs("r_bookings","Booking")
	bookingPrice := g.GenerateModelAs("r_booking_prices","BookingPrice")
	bookingSlot := g.GenerateModelAs("r_booking_slots","BookingSlot")
	bookingEvent := g.GenerateModelAs("r_booking_events","BookingEvent")
	event := g.GenerateModelAs("r_events","EventBooking")

	return []interface{}{
		court,
		courtRate,
		booking,
		bookingPrice,
		bookingSlot,
		bookingEvent,
		event,
	}
}
