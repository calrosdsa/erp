package booking_rest

type BookingPaths struct {
	Base             string
	Detail           string
	Validate         string
	UpdateStatus     string
	UpdatePaidAmount string
	Reschedule       string
	UpdateBookingBatch string 
}

func NewBookingPaths(base string) BookingPaths {
	return BookingPaths{
		Base:             base,
		Detail:           base + "/detail/{id}",
		Validate:         base + "/validate",
		UpdateStatus:     base + "/update-status",
		UpdatePaidAmount: base + "/paid-amount",
		Reschedule:       base + "/reschedule",
		UpdateBookingBatch: base + "/update-booking-batch",
	}
}

type BookingSlotPaths struct {
	Base string
}

func NewBookingSlotPaths(base string) BookingSlotPaths {
	return BookingSlotPaths{
		Base: base,
	}
}

// type CourtRatePaths struct {
// 	Base string
// 	Court string
// }

// func NewCourtRatePaths(base string) CourtRatePaths {
// 	return CourtRatePaths{
// 		Base: base,
// 		Court: base + "/{id}",
// 	}
// }
