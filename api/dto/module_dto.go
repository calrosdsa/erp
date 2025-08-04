package dto

import "erp/gen/db/model"

type (
	ModulesRequest struct {
		DefaultListParams
		Label       string `query:"label"`
		WorkSpaceID string `query:"workspace_id"`
	}
	ModuleDataRequest struct {
		Body ModuleData
	}

	ModuleSearchRequest struct {
		Size         string `query:"size"`
		Query        string `query:"query"`
		LoadModules  bool   `query:"load_modules"`
		LoadEntities bool   `query:"load_entities"`
	}

	ModuleData struct {
		ID       int64               `json:"id" required:"false"`
		Fields   ModuleFields        `json:"fields"`
		Sections []ModuleSectionData `json:"sections"`
	}

	ModuleFields struct {
		Label           string  `json:"label"`
		IconCode        *string `json:"icon_code"`
		IconName        *string `json:"icon_name"`
		Href            string  `json:"href"`
		HasDirectAccess bool    `json:"has_direct_access"`
		Priority        int32   `json:"priority"`
	}

	ModuleDto struct {
		ID              int64   `json:"id"`
		UUID            string  `json:"uuid"`
		Href            string  `json:"href"`
		Label           string  `json:"label"`
		IconCode        *string `json:"icon_code"`
		IconName        *string `json:"icon_name"`
		Status          string  `json:"status"`
		HasDirectAccess bool    `json:"has_direct_access"`
		Priority        int32   `json:"priority"`
	}

	ModuleSectionData struct {
		ModuleID   int64  `json:"module_id"`
		EntityID   int32  `json:"entity_id"`
		EntityName string `json:"entity_name"`
		Name       string `json:"name"`
	}

	ModuleDetailDto struct {
		Module   ModuleDto          `json:"module"`
		Sections []ModuleSectionDto `json:"sections"`
	}

	ModuleSectionDto struct {
		SectionName string `json:"section_name"`
		EntitySectionDto
	}

	EntitySectionDto struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Href string `json:"href"`
	}
)

func ModuleDtoFromModel(m model.Module) ModuleDto {
	return ModuleDto{
		ID:    m.ID,
		Href:  m.Href,
		Label: m.Label,
	}
}
