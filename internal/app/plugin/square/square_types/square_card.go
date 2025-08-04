package squaretypes

import "time"

type SquareCardResponse struct {
	Card struct {
		ID             string `json:"id"`
		CardBrand      string `json:"card_brand"`
		Last4          string `json:"last_4"`
		ExpMonth       int    `json:"exp_month"`
		ExpYear        int    `json:"exp_year"`
		BillingAddress struct {
			PostalCode string `json:"postal_code"`
		} `json:"billing_address"`
		Fingerprint string    `json:"fingerprint"`
		CustomerID  string    `json:"customer_id"`
		MerchantID  string    `json:"merchant_id"`
		Enabled     bool      `json:"enabled"`
		CardType    string    `json:"card_type"`
		PrepaidType string    `json:"prepaid_type"`
		Bin         string    `json:"bin"`
		CreatedAt   time.Time `json:"created_at"`
		Version     int       `json:"version"`
	} `json:"card"`
}