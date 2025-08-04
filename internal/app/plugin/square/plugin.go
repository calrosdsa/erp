package square

import (
	"erp/internal/app/config"
	"erp/internal/app/connection"
	// handlersquare "erp/internal/app/plugin/square/api_square/handler_square"
	// squareservice "erp/internal/app/plugin/square/square_service"
	// squarestrategies "erp/internal/app/plugin/square/square_strategies"
	"erp/internal/app/service/helpers"

	"github.com/asaskevich/EventBus"
)

func InitHandler(
	conn *connection.Connection,
	configService *config.ConfigService,
	bus *EventBus.Bus,
	helpers *helpers.Helpers,
	authMiddleware config.AppMiddleware,
	validateCompanyM config.AppMiddleware,
) {
	// timeout := configService.GetTimeoutAPICall()
	// squareService := squareservice.NewSquareService(conn, configService, &timeout)
	// squareSubscriptionService := squareservice.NewSquareSubscriptionService(conn, configService, &timeout,squareService)
	// apiOptions := configService.GetApiOptions()
	// api := *apiOptions.Api
	// handlersquare.NewHandler(&api, squareService, "/square", "Square",helpers)
	// handlersquare.NewSquareSubscriptionHandler(&api,helpers, squareSubscriptionService,
	// "/square/subscription", "Square Subscription",authMiddleware,validateCompanyM)
}

func Init(
	conn *connection.Connection,
	configService *config.ConfigService,
	helpers *helpers.Helpers,
) *config.PluginApp {
	// timeout := configService.GetTimeoutAPICall()

	// _ = squareservice.NewSquareService(conn, configService, &timeout)


	return &config.PluginApp{
		Name: config.PLUGIN_SQUARE,
		Strategies: config.Strategies{
			// ItemStrategy: squarestrategies.NewSquareItemStrategy(
			// 	conn,
			// 	configService,
			// 	squareService,
			// 	helpers,
			// ),
			// ClientStrategy: squarestrategies.NewSquareClientStrategy(conn, squareService,&timeout),
			// SalesOrderStrategy: squarestrategies.NewSquareSalesOrderStrategy(conn,squareService,timeout,helpers),
			// ItemPriceStrategy: squarestrategies.NewSquareItemPriceStrategy(conn,squareService,&timeout,helpers),
		},
	}
}
