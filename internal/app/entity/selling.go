package entity

import (
	"erp/gen/db/model"
	"time"
)

type SalesOrder struct {
	Base
	Code             string
	PartyID         uint
	Party           model.Party
	OrderType        SalesOrderType
	DeliveryDate     time.Time
	CompanyID        uint
	Company          model.Company
	SalesOrderPlugin []SalesOrderPlugin

	ShippingAddressID *uint
	ShippingAddress   PartyAddress `gorm:"foreignKey:ShippingAddressID,PartyID;references:AddressID,PartyID"`
	BillingAddressID  *uint
	BillingAddress    PartyAddress `gorm:"foreignKey:BillingAddressID,PartyID;references:AddressID,PartyID"`

	Data string `gorm:"-"`
}

type SalesItemLine struct {
	BaseWithoutUuid
	// ItemPrice         ItemPrice
	ItemPriceID       uint
	Rate              int
	ItemQuantity      uint
	Currency          string
	SalesOrderID      uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	
}

type SalesInvoice struct {
	Base
	PartyID       uint
	Party         model.Party
	PaymentDueDate time.Time
	SalesOrderID   uint
	SalesOrder     SalesOrder
	CompanyID      uint
	Company        model.Company
}

type SalesOrderPlugin struct {
	Plugin       string `gorm:"primaryKey;"`
	SalesOrderID uint   `gorm:"primaryKey;" required:"true"`
	Data         string `gorm:"-"`
}

type SalesOrderType string

const (
	ORDER_TYPE_SERVICE  = "SERVICE"  //Service Order: A request for services, often used in maintenance or field service industries.
	ORDER_TYPE_PURCHASE = "PURCHASE" // Purchase Order (PO):A formal request to purchase goods or services from a vendor.

)
