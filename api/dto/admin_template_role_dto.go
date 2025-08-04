package dto

import "time"

type (
	CreateRoleTemplateRequest struct {
		Body struct {
			Name string `json:"name"`
		}
	}

	RoleTemplateDto struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}
)
