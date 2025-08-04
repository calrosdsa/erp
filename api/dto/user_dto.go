package dto

import (
	"erp/gen/db/model"
	"time"
)

type CreateUserRequest struct {
	AuthParams
	Body struct {
		Email      string  `json:"email" required:"true" format:"email"`
		RoleUUID   string  `json:"role_uuid" required:"true"`
		PartyCode  string  `json:"party_code" required:"true" minLength:"1" maxLength:"50"`
		CompanyIds []int64 `json:"company_ids" required:"false"`
		ProfileEditableFields
	}
}

type ProfileEditableFields struct {
	GivenName  string `json:"given_name" required:"true" minLength:"1" maxLength:"50"`
	FamilyName string `json:"family_name" required:"true" minLength:"1" maxLength:"50"`
	// CountryCode string `json:"countryCode" required:"false"`
	PhoneNumber  string         `json:"phone_number" required:"false" maxLength:"50"`
	KeyValueData []KeyValueData `json:"key_value_data" required:"false"`
}

type KeyValueData struct {
	Key   string `json:"key" minLength:"1" maxLength:"50"`
	Value string `json:"value" minLength:"1" maxLength:"50"`
}

type UserDto struct {
	ID         int64      `json:"id"`
	Uuid       string     `json:"uuid"`
	Identifier string     `json:"identifier"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time  `json:"created_at"`
}

func UserDTOFromModel(m *model.User) UserDto {
	// d.Identifier = m.Identifier
	// d.Uuid = m.UUID
	return UserDto{
		Identifier: m.Identifier,
		Uuid:       m.UUID,
		LastLogin:  m.LastLogin,
	}
}
