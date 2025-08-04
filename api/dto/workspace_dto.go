package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	WorkSpaceRequest struct {
		DefaultListParams
		Name string `json:"name"`
	}

	WorkSpaceRequestData struct {
		Body WorkSpaceData
	}
	WorkSpaceData struct {
		ID     int64           `json:"id" required:"false"`
		Fields WorkSpaceFields `json:"fields"`	
		Modules []int64 `json:"modules"`
	}
	WorkSpaceFields struct {
		Name string `json:"name"`
	}

	WorkSpaceDto struct {
		ID        int64           `json:"id"`
		Name      string          `json:"name"`
		CreatedAt time.Time       `json:"created_at"`
		Status    string          `json:"status"`
		Modules   []ModuleDto `json:"modules" gorm:"-"`
	}
)

func WorkspaceFromModel(m model.Workspace) WorkSpaceDto {
	return WorkSpaceDto{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Status:    m.Status,
	}
}
