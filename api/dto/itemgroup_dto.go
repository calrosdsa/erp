package dto

type CreateItemGroupRequest struct {
	AuthParams
	Body struct {
		Name string `json:"name" required:"true" minLength:"3"`
	}
}


