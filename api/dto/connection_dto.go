package dto

type (
	ConnectionDto struct {
		ID             int    `json:"id"`
		EntityID       int64  `json:"entity_id"`
		EntityName     string `json:"entity_name"`
		EntityHref     string `json:"entity_href"`
		EntityHasModal bool   `json:"entity_has_modal"`
		SectionName    string `json:"section_name"`
		Count          int    `json:"count"`
	}

	// ConnectionSectionDto struct {
	// 	SectionName string          `json:"section_name"`
	// 	Connections []ConnectionDto `json:"connections"`
	// }

	// ConnectionsDto struct {
	// 	Sections []ConnectionSectionDto `json:"sections"`
	// }
)
