package entity

import "erp/gen/db/model"

type Address struct {
	Base
	Name                 string
	FullName             string `gorm:"not null;default:null"`
	Company              string
	StreetLine1          string
	StreetLine2          string
	City                 string
	Province             string
	Country              string
	PostalCode           string
	CountryCode          string
	IdentificationNumber string
	PhoneNumber          string
	EmailAddress         string
}

type PartyAddress struct {
	PartyID           uint `gorm:"primaryKey"`
	Party             model.Party
	AddressID         uint `gorm:"primaryKey"`
	Address           Address
	IsShippingAddress bool
	IsBillingAddress  bool
	IsActive          bool
	BaseWithoutID
}
