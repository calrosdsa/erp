package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	// RequestEvents struct {
	// 	PaginationParams
	// 	OptionalQueryParams
	// }

	DeleteEventBatchRequest struct {
		Body struct {
			EventIds []int64 `json:"event_ids"`
			// TargetState string  `json:"target_state" minLength:"1"`
		}
	}
	EventBookingDataRequest struct {
		Body EventBookingData
	}

	EventBookingData struct {
		ID     int64              `json:"id" required:"false"`
		Fields EventBookingFields `json:"fields"`
	}

	EventBookingFields struct {
		Name        string  `json:"name" minLength:"1" maxLength:"255" required:"true"`
		Description *string `json:"description" minLength:"1" maxLength:"255" required:"false"`
	}

	

	EventBookingDto struct {
		ID          int64     `json:"id"`
		UUID        string    `json:"uuid"`
		Status      string    `json:"status"`
		Name        string    `json:"name"`
		Description *string   `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
	}

	EventBookingDetail struct {
		BookingInfo  EventBookingInfo `json:"booking_info"`
		EventBooking EventBookingDto  `json:"event"`
	}
	EventBookingInfo struct {
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		TotalPrice    int    `json:"total_price"`
		TotalPaid     int    `json:"total_paid"`
		TotalDiscount int    `json:"total_discount"`
	}
)

func EventBookingDtoFromModel(m *model.EventBooking) EventBookingDto {
	return EventBookingDto{
		ID:          m.ID,
		UUID:        m.UUID,
		Status:      m.Status,
		Name:        m.Name,
		Description: m.Description,
	}
}
