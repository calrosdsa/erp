package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CreateGroupRequest struct {
		AuthParams
		Body struct {
			GroupData
		}
	}

	EditGroupRequest struct {
		Body struct {
			ID int64 `json:"id" required:"true"`
			GroupData
		}
	}
	GroupData struct {
		Name          string `json:"name" required:"true" minLength:"1" maxLength:"50"`
		// IsGroup       bool   `json:"is_group" required:"true"`
		PartyTypeCode string `json:"party_type_code" required:"true"`
		// ParentID      *int64 `json:"parent_id" required:"false"`
	}

	

	GroupHierarchyDto struct {
		Uuid       string  `json:"uuid"`
		ParentUuid *string `json:"parent_uuid"`
		Name       string  `json:"name"`
		IsGroup    bool    `json:"is_group"`
		Depth      int     `json:"depth"`
	}

	GroupDto struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		IsGroup   bool      `json:"is_group"`
		Ordinal   int       `json:"ordinal"`
		CreatedAt time.Time `json:"created_at"`
		Enabled   bool      `json:"enabled"`
		Uuid      string    `json:"uuid"`
	}
)

func (d *GroupDto) FromModel(c *model.Group) {
	d.ID = c.ID
	d.Name = c.Name
	d.IsGroup = c.IsGroup
	d.Uuid = c.UUID
	d.CreatedAt = c.CreatedAt
	d.Enabled = c.Enabled
	d.Ordinal = int(c.Ordinal)
}

func GroupDtoFromModel(c *model.Group) GroupDto {
	r := GroupDto{}
	r.ID = c.ID
	r.Name = c.Name
	r.IsGroup = c.IsGroup
	r.Uuid = c.UUID
	r.CreatedAt = c.CreatedAt
	r.Enabled = c.Enabled
	r.Ordinal = int(c.Ordinal)
	return r
}
