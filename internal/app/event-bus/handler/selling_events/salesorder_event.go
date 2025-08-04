package sellingevents

import "erp/internal/app/service/services"


type OrderEvent struct {
}

const (
	SALESORDER_EVENT = "salesorder:create"
)

func NewOrderEvent(services *services.Services) {
}
