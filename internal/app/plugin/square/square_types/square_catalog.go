package squaretypes

import "time"

type RetrieveObjectRequest struct {
	Object struct {
		ObjectInfo
		SubscriptionPlanVariationData struct {
			Name               string  `json:"name"`
			Phases             []Phase `json:"phases"`
			SubscriptionPlanID string  `json:"subscription_plan_id"`
		} `json:"subscription_plan_variation_data"`
	} `json:"object"`
}

type ObjectInfo struct {
	Type                  string    `json:"type"`
	ID                    string    `json:"id"`
	UpdatedAt             time.Time `json:"updated_at"`
	CreatedAt             time.Time `json:"created_at"`
	Version               int64     `json:"version"`
	IsDeleted             bool      `json:"is_deleted"`
	PresentAtAllLocations bool      `json:"present_at_all_locations"`
}

type Phase struct {
	UID     string `json:"uid"`
	Cadence string `json:"cadence"`
	Periods int    `json:"periods"`
	Ordinal int    `json:"ordinal"`
	Pricing struct {
		Type       string `json:"type"`
		Price      Amount `json:"price"`
		PriceMoney Amount `json:"price_money"`
	} `json:"pricing"`
}

type Amount struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type SquareTypeObject string

const (
	SQUARE_TYPE_OBJECT       SquareTypeObject = "OBJECT"
	SQUARE_TYPE_SUBSCRIPTION SquareTypeObject = "SUBSCRIPTION"
)
