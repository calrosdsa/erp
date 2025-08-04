package dto

import (
	"erp/gen/db/model"
	"time"
)

//	 CreateRole struct {
//		Body struct {
//			Body
//		}
//	}
type (
	RoleRequestData struct {
		Body RoleData
	}

	RoleData struct {
		ID     int64      `json:"id" required:"false"`
		Fields RoleFields `json:"fields"`
	}

	RoleFields struct {
		Code        string  `json:"code" required:"true"`
		Description *string `json:"description" required:"false"`
		WorkspaceID *int64  `json:"workspace_id" required:"false"`
	}

	EditRolePermissionActions struct {
		AuthParams
		Body struct {
			// RoleID          int64             `json:"roleId"`
			RoleUUID        string           `json:"role_uuid"`
			ActionSelecteds []ActionSelected `json:"actionSelecteds"`
			EntityName      string           `json:"entityName"`
			EntityActions   EntityActionsDto `json:"entity_actions"`
		}
	}

	ActionSelected struct {
		ActionName string `json:"actionName"`
		ActionID   int64  `json:"actionId"`
		Selected   bool   `json:"selected"`
	}

	EntityActionsDto struct {
		Entity  EntityDto   `json:"entity"`
		Actions []ActionDto `json:"actions"`
	}

	RoleDto struct {
		ID          int64      `json:"id"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at"`
		UUID        string     `json:"uuid"`
		Code        string     `json:"code"`
		Description *string    `json:"description"`

		Workspace   *string `json:"workspace"`
		WorkspaceID *int64  `json:"workspace_id"`
	}

	ActionDto struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		EntityID int64  `json:"entity_id"`
	}

	RoleActionDto struct {
		ActionID  int64     `json:"action_id"`
		ActionDto ActionDto `json:"action"`
		RoleId    int64     `json:"role_id"`
	}
)

func ActionDtoFromModel(m *model.Action) ActionDto {
	r := ActionDto{
		ID:       m.ID,
		Name:     m.Name,
		EntityID: m.EntityID,
	}
	return r
}

func RoleActionDTOFromModel(m *model.RoleAction) RoleActionDto {
	r := RoleActionDto{}
	r.ActionID = m.ActionID
	r.RoleId = m.RoleID
	if m.Action.ID != 0 {
		r.ActionDto = ActionDtoFromModel(&m.Action)
	}
	return r
}

func RoleDTOFromModel(m *model.Role) RoleDto {
	p := RoleDto{}
	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.Code = m.Code
	p.UUID = m.UUID
	p.Description = m.Description
	p.UpdatedAt = m.UpdatedAt
	p.WorkspaceID = m.WorkspaceID
	return p
}
