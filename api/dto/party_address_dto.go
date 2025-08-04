package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	AddressDataRequest struct {
		Body AddressData
	}

	AddressData struct {
		ID          int64         `json:"id" required:"false"`
		ReferenceID *int64        `json:"reference_id" required:"false"`
		Action      string        `json:"action" required:"false"`
		Fields      AddressFields `json:"fields"`
	}

	AddressFields struct {
		Title                string  `json:"title"`
		IsShippingAddress    bool    `json:"is_shipping_address"`
		IsBillingAddress     bool    `json:"is_billing_address"`
		StreetLine1          string  `json:"street_line1"`
		StreetLine2          string  `json:"street_line2"`
		City                 string  `json:"city"`
		Company              *string `json:"company" required:"false"`
		Province             *string `json:"province" required:"false"`
		PostalCode           *string `json:"postal_code" required:"false"`
		PhoneNumber          *string `json:"phone_number" required:"false"`
		CountryCode          *string `json:"country_code" required:"false"`
		IdentificationNumber *string `json:"identification_number" required:"false"`
		Email                *string `json:"email" required:"false"`
	}

	AddressRequestData struct {
		Title                string `json:"title" required:"true" minLength:"1" maxLength:"50"`
		Company              string `json:"company" required:"false"`
		StreetLine1          string `json:"street_line_1" required:"true" minLength:"1" maxLength:"50"`
		StreetLine2          string `json:"street_line_2" required:"true" minLength:"1" maxLength:"50"`
		City                 string `json:"city" required:"true" minLength:"1" maxLength:"50"`
		Province             string `json:"province" required:"false"`
		CountryCode          string `json:"country_code" required:"false"`
		PostalCode           string `json:"postal_code" required:"false"`
		PhoneNumber          string `json:"phone_number" required:"false"`
		Email                string `json:"email" required:"false"`
		IdentificationNumber string `json:"identification_number" required:"false"`
		IsBillingAddress     bool   `json:"is_billing_address" required:"true"`
		IsShippingAddress    bool   `json:"is_shipping_address" required:"true"`
		Enabled              bool   `json:"enabled" required:"true"`
	}

	BillingAddressDto struct {
		IsShippingAddress bool      `json:"is_shipping_address"`
		IsBillingAddress  bool      `json:"is_billing_address"`
		IsActive          bool      `json:"is_active"`
		CreatedAt         time.Time `json:"created_at"`

		Company              *string `json:"company"`
		StreetLine1          string  `json:"street_line1"`
		StreetLine2          string  `json:"street_line2"`
		City                 string  `json:"city"`
		Province             *string `json:"province"`
		PostalCode           *string `json:"postal_code"`
		PhoneNumber          *string `json:"phone_number"`
		CountryCode          *string `json:"country_code"`
		IdentificationNumber *string `json:"identification_number"`
		Email                *string `json:"email"`
		// Address           AddresDto `json:"address"`
	}

	AddressDto struct {
		ID                   int64   `json:"id"`
		UUID                 string  `json:"uuid"`
		Title                string  `json:"title"`
		City                 string  `json:"city"`
		StreetLine1          string  `json:"street_line1"`
		StreetLine2          string  `json:"street_line2"`
		IsShippingAddress    bool    `json:"is_shipping_address"`
		IsBillingAddress     bool    `json:"is_billing_address"`
		Company              *string `json:"company"`
		Province             *string `json:"province"`
		PostalCode           *string `json:"postal_code"`
		PhoneNumber          *string `json:"phone_number"`
		CountryCode          *string `json:"country_code"`
		IdentificationNumber *string `json:"identification_number"`
		Email                *string `json:"email"`
		Status               string  `json:"status"`
	}
)

func AddressDtoFromModel(m *model.Address) AddressDto {
	r := AddressDto{}
	r.ID = m.ID
	r.UUID = m.UUID
	r.Title = m.Title
	r.Company = m.Company
	r.StreetLine1 = m.StreetLine1
	r.StreetLine2 = m.StreetLine2
	r.City = m.City
	r.Province = m.Province
	r.PostalCode = m.PostalCode
	r.PhoneNumber = m.PhoneNumber
	r.CountryCode = m.CountryCode
	r.IdentificationNumber = m.IdentificationNumber
	r.IsShippingAddress = m.IsShippingAddress
	r.IsBillingAddress = m.IsBillingAddress
	r.Email = m.Email
	r.Status = m.Status
	return r
}
