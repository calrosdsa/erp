package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CourtsRequest struct {
		DefaultListParams
		Name      string `query:"name" required:"false"`
		CreatedAt string `query:"created_at" required:"false"`
		UpdatedAt string `query:"updated_at" required:"false"`
	}
	CreateCourtRequest struct {
		Body CreateCourtBody
	}
	CreateCourtBody struct {
		Name    string `json:"name"`
	}
	EditCourtRequest struct {
		Body EditCourtBody
	}
	EditCourtBody struct {
		Name    string `json:"name" required:"true" minLength:"1" maxLength:"50"`
		CourtID int64  `json:"court_id" required:"true"`
	}

	CourtDto struct {
		ID        int64     `json:"id"`
		UUID      string    `json:"uuid"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		Status   string      `json:"status"`
	}

	AvailableCourtDto struct {
		CourtDto
		TotalPrice int `json:"total_price"`
	}
)

func CourtDtoFromModel(m *model.Court) CourtDto {
	return CourtDto{
		ID:        m.ID,
		UUID:      m.UUID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Status:   m.Status,
	}
}
