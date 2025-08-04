package buying

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
	// NewSupplierHandler(api, services, helpers, "/supplier", []string{"Supplier", "Buying"}, middlewares)
	NewPurchaseHandler(api, services, helpers, "/purchase/order/", []string{"Purchase Order", "Buying"}, middlewares)
}
