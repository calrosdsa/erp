package party

import (
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"

	"github.com/danielgtaylor/huma/v2"
)

func NewHandlers(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
) {
	// NewPartyAddressHandler(api, services, helpers, "/party/address", []string{"Address", "Party"}, middlewares)
	// NewPartyContactHandler(api, services, helpers, "/party/contact", []string{"Contact", "Party"}, middlewares)
	// NewPartyHandler(api, services, helpers, "/party", []string{"Party"}, middlewares)
	// NewPurchaseHandler(api, services, helpers, "/purchase/order/", []string{"Purchase Order", "Buying"}, middlewares)
}
