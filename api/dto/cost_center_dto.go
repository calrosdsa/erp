package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CreateCostCenterRequet struct {
		Body struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		}
	}
	CostCenterDto struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}


)

func CostCenterDtoFromModel(m *model.CostCenter) CostCenterDto {
	return CostCenterDto{
		ID:        m.ID,
		Name:      m.Name,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}
}
