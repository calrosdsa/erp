package dto

import (
	"erp/api/resolvers"
	"erp/internal/app/entity"

	"gorm.io/gorm"
)

type CreateClientRequest struct {
	Body struct {
		ClientRequestDto ClientRequestDto `json:"client"`
	}
}

type ClientRequestDto struct {
	GivenName    string  `json:"givenName"`
	FamilyName   string  `json:"familyName"`
	CompanyName  string  `json:"companyName"`
	EmailAddress string  `json:"email"`
	PhoneNumber  string  `json:"phoneNumber"`
	Country      Country `json:"country" required:"false"`

	DeleteAt gorm.DeletedAt `json:"deleteAt" required:"false"`

	Plugins []entity.CompanyPlugins `json:"plugins" required:"false"`

	KeyValues []entity.ClientKeyValueData `json:"keyValues" required:"false"`
	Metadata  string                      `json:"metadata" required:"false"`
}

type EditClientRequest struct {
	AuthParams
	Body struct {
		ClientEditableFields
	}
}

type ClientEditableFields struct {
	GivenName        string                `json:"givenName"`
	FamilyName       string                `json:"familyName"`
	OrganizationName string                `json:"organizationName"`
	PhoneNumber      resolvers.PhoneNumber `json:"phoneNumber"`
}

type Country struct {
	Code  string `json:"code"`
	Label string `json:"label" required:"false"`
	Phone string `json:"phone" required:"false"`
}
