package squaretypes

import (
	"erp/internal/app/entity"
	"time"
)

type SquareCustomerMetadata struct {
	PlanVariationId string `json:"objectId"`
	ItemGroupUuid   string `json:"itemGroupUuid"`
	Type            SquareTypeObject `json:"type"`
	CardRequest     struct {
		LocationId        string `json:"locationId"`
		SourceId          string `json:"sourceId"`
		VerificationToken string `json:"verificationToken"`
		IdEmpotencyKey    string `json:"idempotencyKey"`
	} `json:"cardRequest"`
}

type SquareCustomerResponse struct {
	Customer struct {
		ID           string    `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		GivenName    string    `json:"given_name"`
		FamilyName   string    `json:"family_name"`
		EmailAddress string    `json:"email_address"`
		PhoneNumber  string    `json:"phone_number"`
		CompanyName  string    `json:"company_name"`
		Preferences  struct {
			EmailUnsubscribed bool `json:"email_unsubscribed"`
		} `json:"preferences"`
		CreationSource string `json:"creation_source"`
		Version        int    `json:"version"`
	} `json:"customer"`
}

type SquareCustomerRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	GivenName      string `json:"given_name"`
	FamilyName     string `json:"family_name"`
	CompanyName    string `json:"company_name"`
	EmailAddress   string `json:"email_address"`
	// PhoneNumber    string `json:"phone_number"`
}

func (c *SquareCustomerRequest) FromEntity(e *entity.Client) {
	c.CompanyName = e.OrganizationName
	c.EmailAddress = e.EmailAddress
	// c.PhoneNumber = e.PhoneNumber
	c.FamilyName = e.FamilyName
	c.GivenName = e.GivenName
	c.IdempotencyKey = e.Uuid
}
