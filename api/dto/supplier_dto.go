package dto

import (
	"erp/gen/db/model"
)

type (
	SupplierDataRequest struct {
		AuthParams
		Body SupplierData
	}

	SupplierData struct {
		ID       int64          `json:"id" required:"false"`
		Fields   SupplierFields `json:"fields" required:"true"`
		Contacts []ContactData  `json:"contacts"`
	}

	SupplierFields struct {
		Name    string `json:"name" required:"true" minLength:"1" maxLength:"50"`
		GroupID *int64 `json:"group_id" required:"false"`
	}

	// CustomerDataRequest struct {
	// 	AuthParams
	// 	Body CustomerData
	// }

	// CustomerData struct {
	// 	ID       int64          `json:"id" required:"false"`
	// 	Fields   CustomerFields `json:"fields"`
	// 	Contacts []ContactData  `json:"contacts"`
	// }
	// // CustomerData struct {}

	// CustomerFields struct {
	// 	Name         string `json:"name" minLength:"1" maxLength:"50" required:"true"`
	// 	CustomerType string `json:"customer_type" minLength:"1" maxLength:"50" required:"true"`
	// 	GroupID      *int64 `json:"group_id" required:"false"`
	// }

	SupplierDto struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Uuid      string `json:"uuid"`
		Status    string `json:"status"`
		CreatedAt string `json:"created_at"`

		Email       *string `json:"email"`
		PhoneNumber *string `json:"phone_number"`

		Group     string `json:"group"`
		GroupID   int64  `json:"group_id"`
		GroupUUID string `json:"group_uuid"`
	}
)

func (d *SupplierDto) FromModel(c *model.Supplier) *SupplierDto {
	d.Name = c.Name
	d.Uuid = c.UUID
	return d
}

func SupplierDtoFromModel(c *model.Supplier) SupplierDto {
	r := SupplierDto{}
	r.Name = c.Name
	r.ID = c.ID
	r.Uuid = c.UUID
	return r
}

// func (d *SupplierDto) SetGroup(c GroupDto) *SupplierDto {
// 	d.Group = c
// 	return d
// }
