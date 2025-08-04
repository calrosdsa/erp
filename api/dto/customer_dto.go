package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	CustomerDataRequest struct {
		AuthParams
		Body CustomerData
	}

	CustomerData struct {
		ID       int64          `json:"id" required:"false"`
		Fields   CustomerFields `json:"fields"`
		Contacts []ContactData  `json:"contacts"`
	}
	// CustomerData struct {}

	CustomerFields struct {
		Name         string `json:"name" minLength:"1" maxLength:"50" required:"true"`
		CustomerType string `json:"customer_type" minLength:"1" maxLength:"50" required:"true"`
		GroupID      *int64 `json:"group_id" required:"false"`
	}

	EditCustomerRequest struct {
		Body struct {
			ID int64 `json:"customer" required:"true"`
			CustomerFields
		}
	}
	// EditCustomerBody struct {
	// 	Name         string `json:"name" required:"true"`
	// 	CustomerType string `json:"customer_type" minLength:"1" maxLength:"50" required:"true"`
	// 	GroupID      *int64 `json:"group_id" required:"false"`
	// }

	CustomerType struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	CustomerDto struct {
		ID           int64     `json:"id"`
		UUID         string    `json:"uuid"`
		Name         string    `json:"name"`
		CreatedAt    time.Time `json:"created_at"`
		CustomerType string    `json:"customer_type"`
		Status       string    `json:"status"`
		//Group
		GroupID   *int64  `json:"group_id"`
		GroupUUID *string `json:"group_uuid" required:"false"`
		GroupName *string `json:"group_name" required:"false"`
	}
)

func CustomerDtoFromModel(m *model.Customer) CustomerDto {
	r := CustomerDto{}
	r.ID = m.ID
	r.UUID = m.UUID
	r.Name = m.Name
	r.CreatedAt = m.CreatedAt
	r.CustomerType = m.CustomerType

	return r
}
