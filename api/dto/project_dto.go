package dto

import "erp/gen/db/model"

type (
	CreateProjectRequest struct {
		Body struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		}
	}
	ProjectDto struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
)

func ProjectDtoFromModel(m *model.Project) ProjectDto{
	return ProjectDto{
		ID:m.ID,
		Name:m.Name,
		Status: m.Status,
	}
}
