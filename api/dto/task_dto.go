package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	TasksRequest struct {
		DefaultListParams
		TaskFilterParams
		CreatedAt string `query:"created_at" required:"false"`
		UpdatedAt string `query:"updated_at" required:"false"`
	}

	TaskDataRequest struct {
		Body TaskData
	}

	TaskData struct {
		ID     int64      `json:"id"`
		Fields TaskFields `json:"fields"`
	}

	TaskFields struct {
		ProjectID   int64      `json:"project_id" required:"true"`
		Assignee    *int64     `json:"assignee" required:"false"`
		Title       string     `json:"title" required:"true"`
		Description *string    `json:"description" required:"false"`
		StageID     int32      `json:"stage_id" required:"false"` // Stage ID for database
		Priority    *string    `json:"priority" required:"false"`
		DueDate     *time.Time `json:"due_date" required:"false"`
		Index       int32      `json:"index" required:"false"` // Task index within stage
	}

	TaskDetailDto struct {
		Task TaskDto `json:"task"`
	}

	TaskDto struct {
		ID   int64  `json:"id"`
		UUID string `json:"uuid"`

		ProjectID int64  `json:"project_id"`
		Project   string `json:"project"`

		AssigneeID *int64 `json:"assignee_id"`
		AssigneeGivenName string `json:"assignee_given_name"`
		AssigneeFamilyName string `json:"assignee_family_name"`

		Title       string     `json:"title"`
		Description *string    `json:"description"`
		Priority    *string    `json:"priority"`
		DueDate     *time.Time `json:"due_date"`

		StageID    int32  `json:"stage_id"`
		StageIndex int32  `json:"stage_index"`
		Stage      string `json:"stage"`

		Index int32 `json:"index"`

		CompanyID int64      `json:"company_id"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
	}

	// Kanban-specific DTOs
	KanbanColumnDto struct {
		Count int64     `json:"count"`
		Tasks []TaskDto `json:"tasks"`
	}

	KanbanBoardDto struct {
		ProjectID int64             `json:"project_id"`
		Columns   []KanbanColumnDto `json:"columns"`
		Total     int64             `json:"total"`
	}

	// Task filtering DTOs
	TaskFilterParams struct {
		Status      string `query:"status,omitempty" doc:"Filter by task status" required:"false"`
		Assignee    int64  `query:"assignee,omitempty" doc:"Filter by assignee ID" required:"false"`
		Priority    string `query:"priority,omitempty" doc:"Filter by priority" required:"false"`
		Search      string `query:"search,omitempty" doc:"Search in title and description" required:"false"`
		DueDateFrom string `query:"due_date_from,omitempty" doc:"Due date range start (YYYY-MM-DD)" required:"false"`
		DueDateTo   string `query:"due_date_to,omitempty" doc:"Due date range end (YYYY-MM-DD)" required:"false"`
		CreatedFrom string `query:"created_from,omitempty" doc:"Created date range start (YYYY-MM-DD)" required:"false"`
		CreatedTo   string `query:"created_to,omitempty" doc:"Created date range end (YYYY-MM-DD)" required:"false"`
	}
)

func TaskFromModel(m model.Task) TaskDto {
	return TaskDto{
		ID:          m.ID,
		UUID:        m.UUID,
		ProjectID:   m.ProjectID,
		AssigneeID:    m.Assignee,
		Title:       m.Title,
		Description: m.Description,
		Priority:    m.Priority,
		DueDate:     m.DueDate,
		CompanyID:   m.CompanyID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
