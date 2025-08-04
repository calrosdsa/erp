package dto

import (
	"erp/internal/app/entity"
	"time"

	"gorm.io/gorm"
)

type RequestSalesOrderDetail struct {
	AuthParams
	OrderCode string `path:"code"`
}

type ResponseSalesOrderDetail struct {
	Body struct {
		SalesOrder      entity.SalesOrder      `json:"order"`
		SalesOrderItems []entity.SalesItemLine `json:"lines"`
	}
}

type CreateSalesOrderRequest struct {
	AuthParams
	Body struct {
		CreateSalesOrderBody
	}
}

type CreateSalesOrderBody struct {
	ClientID          uint                  `json:"clientId" required:"true"`
	OrderType         entity.SalesOrderType `json:"orderType" required:"true"`
	DeliveryDate      time.Time             `json:"deliveryDate"`
	IsValid           bool                  `json:"isValid" required:"false"`
	SalesItemLines    []SalesItemLineDto    `json:"sales_item_lines"`
	Plugins           []string              `json:"plugins" required:"false"`
	DeleteAt          gorm.DeletedAt        `required:"false"`
	ShippingAddressID *uint                  `json:"shippingAddressId" required:"false"`
	BillingAddressID  *uint                  `json:"billingAddressId" required:"false"`
}

type SalesItemLineDto struct {
	ItemPriceID   uint    `json:"itemPriceId" required:"true"`
	Rate          float64 `json:"rate" required:"true"`
	Currency      string  `json:"currency" required:"true"`
	ItemQuanitity uint    `json:"itemQuantity"`
}
