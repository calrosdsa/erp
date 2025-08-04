package api

import (
	// "erp/internal/api/common"
	"context"
	"erp/api/handlers/accounting"
	"erp/api/handlers/buying"
	"erp/api/handlers/client"
	// "erp/api/handlers/company"
	"erp/api/handlers/domain"
	// "erp/api/handlers/group"
	"erp/api/handlers/integrations/cuatropf"

	// "erp/api/handlers/invoice"
	"erp/gen/proto"
	"erp/pkg/discovery/consul"

	// "erp/api/handlers/integrations/teclumobility"
	// "erp/api/handlers/order"
	// "erp/api/handlers/party"
	_pluginhandler "erp/api/handlers/plugin"
	"erp/api/handlers/selling"
	"erp/api/handlers/stock"
	"erp/api/handlers/uom"
	// "erp/api/handlers/user"
	// "erp/api/handlers/warehouse"
	"erp/api/middlewares"

	// "erp/internal/api/middlewares"
	"erp/internal/app/config"
	"erp/internal/app/connection"
	"erp/internal/app/gateway/accounting/grpc"
	"erp/internal/app/plugin"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/internal/domain/repository"
	"fmt"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/danielgtaylor/huma/v2"
	// "github.com/labstack/echo/v4"
	// echojwt "github.com/labstack/echo-jwt/v4"
)

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func Init(
	services *services.Services,
	ConfigModule *config.ConfigModule,
	connection *connection.Connection,
	helpers *helpers.Helpers,
	plugin *plugin.PluginModule,
	sessionService repository.SessionService,
) {
	apiOptions := ConfigModule.ConfigService.GetApiOptions()

	// e := echo.New()
	// e := apiOptions.EchoServer
	api := apiOptions.Api
	authenticateMiddleware := middlewares.NewMiddlewares(sessionService, api, helpers.Jwt)

	plugin.InitHandlers(authenticateMiddleware.Authenticate, authenticateMiddleware.ValidateActiveCompany)

	huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	// `// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// // })`
	domain.NewDomainHandler(&api, services, helpers, "", []string{"Domain"}, huma.Middlewares{authenticateMiddleware.Authenticate})
	// group.NewGroupHandler(&api, services, helpers, "/group", []string{"Group"}, huma.Middlewares{authenticateMiddleware.Authenticate})
	// order.NewOrderHandler(&api, services, helpers, "/order", []string{"Group"}, huma.Middlewares{authenticateMiddleware.Authenticate})
	// authG.Use(middlewares.Authenticate(connection.Adapter))
	// user.NewProfileHandler(&api, services, helpers, "/user/profile", []string{"Profile"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	// account.NewHandler(&api, services, helpers, "/account", "Account", huma.Middlewares{authenticateMiddleware.Authenticate})
	// account.NewRoleHandler(&api, services, helpers, "/role", "Role", huma.Middlewares{authenticateMiddleware.Authenticate})

	client.NewHandler(&api, services, helpers, "/client", []string{"Client"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	// company.NewHandler(&api, services, helpers, "/company", "Company", huma.Middlewares{authenticateMiddleware.Authenticate})

	_pluginhandler.NewHandler(&api, services, helpers, "/plugin", "Plugin",
		authenticateMiddleware.Authenticate, authenticateMiddleware.ValidateActiveCompany)

	// invoice.NewInvoiceHandler(&api, services, helpers, "/invoice", []string{"Invoice"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	// party.NewHandlers(&api, services, helpers, huma.Middlewares{authenticateMiddleware.Authenticate})
	buying.NewHandlers(&api, services, helpers, huma.Middlewares{authenticateMiddleware.Authenticate})
	// selling.NewSellingHandlers(&api, services, helpers, huma.Middlewares{authenticateMiddleware.Authenticate})

	// stock.NewItemGroupHandler(&api, services, helpers, "/stock/item-group", []string{"Stock", "Item Group"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	// stock.NewItemHandler(&api, services, helpers, "/stock/item", []string{"Stock", "Item"}, huma.Middlewares{authenticateMiddleware.Authenticate})
	// stock.NewPriceListHandler(&api, services, helpers, "/stock/item/price-list", []string{"Stock", "Price List"}, huma.Middlewares{authenticateMiddleware.Authenticate})
	// stock.NewItemPriceHandler(&api, services, helpers, "/stock/item/item-price", []string{"Stock", "Item Price"}, huma.Middlewares{authenticateMiddleware.Authenticate})
	stock.NewItemAttributeHandler(&api, services, helpers, "/stock/item/item-attribute", []string{"Stock", "Item Attribute"}, huma.Middlewares{authenticateMiddleware.Authenticate})
	stock.NewItemVariantHandler(&api, services, helpers, "/stock/item/variant", []string{"Stock", "Item Variant"}, huma.Middlewares{authenticateMiddleware.Authenticate})
	stock.NewItemStockHandler(&api, services, helpers, "/stock/item/level", []string{"Stock", "Item Stock"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	// warehouse.NewWareHouseHandler(&api, services, helpers, "/stock/warehouse", []string{"Warehouse", "Stock"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	selling.NewSaleOrderHandler(&api, services, helpers, "/selling/salesorder", []string{"Selling", "Sales Order"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	accounting.NewTaxHandler(&api, services, helpers, "/accounting/tax", []string{"Tax", "Accounting"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	cuatropf.NewHandler(&api, services, helpers, "/cuatropf", "Cuatropf")
	// teclumobility.NewHandler(&api, services, helpers, "/teclumobity", "Teclu Mobility")

	uom.NewHandler(&api, services, helpers, "/uom", []string{"UOM"}, huma.Middlewares{authenticateMiddleware.Authenticate})

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	transactionGrpcClient := grpc.New(registry)
	_, err = transactionGrpcClient.SaveTransaction(context.Background(), &proto.TransactionLedger{
		LedgerNo:   1,
		LedgerNoDe: 2,
		Amount:     100,
	})
	fmt.Println("ERROR SENDING TRANSACTION", err)

	// if err := e.Start(fmt.Sprintf(":%d", apiOptions.Port)); err != http.ErrServerClosed {
	// 	// log.Fatal(err)
	// 	fmt.Println(err)
	// }
}
