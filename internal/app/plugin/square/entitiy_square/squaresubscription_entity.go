package entitysquare

import "erp/internal/app/entity"

type SquareSubscription struct {
	entity.BaseWithoutUuid
	SubscriptionId string 
	CustomerId string
	PlanVariationId string 
	Status string 
}


type SquareSubscriptionState string

const (
	SQUARE_ACTIVE_SUBSCRIPTION SquareSubscriptionState = "ACTIVE"
)