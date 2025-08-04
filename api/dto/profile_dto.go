package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	ProfilesRequest struct {
		DefaultListParams
		FullName string `query:"full_name"`
	}

	UpdateProfileRequest struct {
		AuthParams
		Body struct {
			ProfileFields EditableProfileFields `json:"profile"`
		}
	}

	EditableProfileFields struct {
		GivenName   string  `json:"given_name" minLength:"1" maxLength:"50" required:"true"`
		FamilyName  string  `json:"family_name"  maxLength:"50" required:"true"`
		PhoneNumber *string `json:"phone_number"  maxLength:"50" required:"false"`
	}

	ProfileL struct {
		ID           int       `json:"id"`
		Uuid         string    `json:"uuid"`
		GivenName    string    `json:"givenName"`
		FamilyName   string    `json:"familyName"`
		EmailAddress string    `json:"emailAddress"`
		PhoneNumber  string    `json:"phoneNumber"`
		PartyCode    string    `json:"partyCode"`
		PartyName    string    `json:"partyName"`
		CreatedAt    time.Time `json:"createdAt"`
	}

	ProfileDto struct {
		ID           int64   `json:"id"`
		Uuid         string  `json:"uuid"`
		GivenName    string  `json:"given_name"`
		FamilyName   string  `json:"family_name"`
		FullName     string  `json:"full_name"`
		EmailAddress string  `json:"email"`
		PhoneNumber  *string `json:"phone_number"`
	}
)

func ProfileDTOFromModel(m *model.Profile) ProfileDto {
	p := ProfileDto{}
	p.ID = m.ID
	p.EmailAddress = m.EmailAddress
	p.FamilyName = m.FamilyName
	p.GivenName = m.GivenName
	p.Uuid = m.UUID
	if m.PhoneNumber != nil {
		p.PhoneNumber = m.PhoneNumber
	}
	return p
}
