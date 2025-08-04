package eventbooking_rest

type EventBookingPaths struct {
	Base   string
	Detail string
	UpdateStatus string 
	DeleteInBatch string 
}

func NewEventBookingPaths(base string) EventBookingPaths {
	return EventBookingPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
		DeleteInBatch: base + "/delete-in-batch",
	}
}
