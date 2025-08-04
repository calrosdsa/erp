package dto

type CuatropfSubscriptionRequest struct {
	CompanyUuid string `path:"companyUuid" required:"true"`

	Body struct {
		ClientRequestDto
	}
}

type TecluMobilityRequestItemPrice struct {
	Body struct {
		Type          string `json:"type"`
		ItemCode      string `json:"itemCode"`
		BillingPeriod uint   `json:"billingPeriod" required:"false"`
		
	}
}

type TecluMobilityOrderRequest struct {
	AcceptLanguage string `header:"Accept-Language"`
	Body struct {
		OrderData        OrderData2        `json:"order"`
		ClientRequestDto ClientRequestDto `json:"client"`
		BillingData      BillingData      `json:"billing" required:"false"`
	}
}

type TecluMobilityOrderResponse struct {
	Body struct {
		PaymentUrl string `json:"paymentUrl"`
	}
}

type OrderData2 struct {
	OrderLine []OrderLineData `json:"orderLine"`
}

type OrderLineData struct {
	ItemPriceId uint `json:"itemPriceId"`
	Quantity    uint `json:"quantity"`
}

type CustomerDataI struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	CompanyName string `json:"companyName"`
	Position    string `json:"cargo"`
}

type BillingData struct {
	Country                 string `json:"country"`
	Address                 string `json:"address"`
	TaxIdentificationNumber string `json:"taxIdNumber"`
	City                    string `json:"city"`
	Estado                  string `json:"estado"`
	PostalCode              string `json:"postalCode"`
}
