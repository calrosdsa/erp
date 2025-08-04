package dto

import (
	"time"
	"erp/gen/db/model"
)

type (
	CreateTaskRequest struct {
		Body struct {
			ProjectID   int64   `json:"project_id"`
			Assignee    *int64  `json:"assignee,omitempty"`
			Title       string  `json:"title"`
			Description *string `json:"description,omitempty"`
			Status      string  `json:"status"`
			Priority    *string `json:"priority,omitempty"`
			DueDate     *string `json:"due_date,omitempty"` // ISO date string
		}
	}
	TaskDto struct {
		ID          int64     `json:"id"`
		ProjectID   int64     `json:"project_id"`
		Assignee    *int64    `json:"assignee,omitempty"`
		Title       string    `json:"title"`
		Description *string   `json:"description,omitempty"`
		Status      string    `json:"status"`
		Priority    *string   `json:"priority,omitempty"`
		DueDate     *string   `json:"due_date,omitempty"` // ISO date string
		CompanyID   int64     `json:"company_id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	}
)

func TaskDtoFromModel(m *model.Task) TaskDto {
	dto := TaskDto{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
		Assignee:    m.Assignee,
		Title:       m.Title,
		Description: m.Description,
		Status:      m.Status,
		Priority:    m.Priority,
		CompanyID:   m.CompanyID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	
	// Convert time.Time to ISO date string for DueDate
	if m.DueDate != nil {
		dueDateStr := m.DueDate.Format("2006-01-02")
		dto.DueDate = &dueDateStr
	}
	
	return dto
}