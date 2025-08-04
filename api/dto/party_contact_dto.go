package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	ContactsRequest struct {
		DefaultListParams
		Name string `query:"name"`
		PartyID string  `query:"party_id"`
	}

	ContactDataRequest struct {
		Body ContactData
	}

	ContactBulkDataRequest struct {
		Body ContactBulkData
	}

	ContactBulkData struct {
		PartyID  int64         `json:"party_id" required:"false"`
		Contacts []ContactData `json:"contacts"`
	}

	ContactData struct {
		ID          int64         `json:"id" required:"false"`
		ReferenceID *int64        `json:"reference_id" required:"false"`
		Action      string        `json:"action" required:"false"`
		Fields      ContactFields `json:"fields"`
	}

	ContactFields struct {
		Name        string  `json:"name" required:"true" minLength:"1" maxLength:"50"`
		Gender      *string `json:"gender" required:"false"`
		Email       *string `json:"email" required:"false" format:"email"`
		PhoneNumber *string `json:"phone_number" required:"false"`
	}

	ContactDto struct {
		ID          int64     `json:"id"`
		UUID        string    `json:"uuid"`
		Name        string    `json:"name"`
		Gender      *string   `json:"gender"`
		Email       *string   `json:"email"`
		PhoneNumber *string   `json:"phone_number"`
		CreatedAt   time.Time `json:"created_at"`
	}
)

func ContactDtoFromModel(m *model.Contact) ContactDto {
	r := ContactDto{}
	r.ID = m.ID
	r.UUID = m.UUID
	r.Email = m.Email
	r.Name = m.Name
	r.Gender = m.Gender
	r.PhoneNumber = m.PhoneNumber
	r.CreatedAt = m.CreatedAt
	return r
}
