package dto

type (
	AddEntityActionRequest struct {
		Body struct {
			Name     string `json:"name"`
			EntityID int64  `json:"entity_id"`
		}
	}

	CreateEntityRequest struct {
		Body struct {
			Name string `json:"name"`
		}
	}

	EntityDto struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Href string `json:"href"`
	}

	EntityDetailDto struct {
		Entity  EntityDto   `json:"entity"`
		Actions []ActionDto `json:"actions"`
	}
)
