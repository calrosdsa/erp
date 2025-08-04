package squaretypes

import (
	entitysquare "erp/internal/app/plugin/square/entitiy_square"
	"time"
)

type SquareSubscriptionResponse struct {
	Subscription struct {
		ID                 string    `json:"id"`
		LocationID         string    `json:"location_id"`
		CustomerID         string    `json:"customer_id"`
		StartDate          string    `json:"start_date"`
		CanceledDate       string    `json:"canceled_date"`
		ChargedThroughDate string    `json:"charged_through_date"`
		Status             string    `json:"status"`
		InvoiceIds         []string  `json:"invoice_ids"`
		Version            int       `json:"version"`
		CreatedAt          time.Time `json:"created_at"`
		CardID             string    `json:"card_id"`
		Timezone           string    `json:"timezone"`
		Source             struct {
			Name string `json:"name"`
		} `json:"source"`
		Actions []SubscriptionActions `json:"actions"`
		MonthlyBillingAnchorDate int    `json:"monthly_billing_anchor_date"`
		PlanVariationID          string `json:"plan_variation_id"`
	} `json:"subscription"`
}

type SubscriptionActions struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	EffectiveDate string `json:"effective_date"`
}

type SalesOrderSquareSubscription struct {
	SquareSubscription entitysquare.SquareSubscription `json:"squareSubscription"`
	Subscription       SquareSubscriptionResponse      `json:"subscription"`
}


