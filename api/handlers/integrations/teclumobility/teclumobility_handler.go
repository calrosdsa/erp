package teclumobility

// import (
// 	"context"
// 	"database/sql"
// 	"erp/api/common"
// 	"erp/api/dto"
// 	"erp/internal/app/config"
// 	"erp/internal/app/entity"
// 	"erp/internal/app/service/helpers"
// 	"erp/internal/app/service/services"
// 	addressservice "erp/internal/app/service/services/address_service"
// 	clientservice "erp/internal/app/service/services/client_service"
// 	"erp/internal/app/service/services/company_service"
// 	integrationservice "erp/internal/app/service/services/integration_service"
// 	sellingservice "erp/internal/app/service/services/selling_service"
// 	"net/http"
// 	"time"

// 	"github.com/danielgtaylor/huma/v2"
// 	"gorm.io/gorm"
// )

// // optional code omitted

// type TecluMobilityHandler struct {
// 	locale               helpers.Locale
// 	tecluMobilityService *integrationservice.TecluMobilityService
// 	companyService       *company_service.CompanyService
// 	clientService        *clientservice.ClientService
// 	orderService         *sellingservice.SalesOrderService
// 	addressService       *addressservice.AddressService
// }

// func NewHandler(
// 	api *huma.API,
// 	services *services.Services,
// 	helpers *helpers.Helpers,
// 	base string,
// 	tag string,
// 	// middlewares huma.Middlewares,
// ) {
// 	paths := NewTecluMobilityPath(base)
// 	handler := TecluMobilityHandler{
// 		tecluMobilityService: services.TecluMobilityService,
// 		locale:               helpers.Locale,
// 		companyService:       services.CompanyService,
// 		clientService:        services.ClientService,
// 		orderService:         services.SalesOrderService,
// 		addressService:       services.AddressService,
// 	}
// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "get-item-price",
// 		Method:        http.MethodPost,
// 		Path:          paths.ItemPrice,
// 		Summary:       "Retrieve item price",
// 		Tags:          []string{tag},
// 		DefaultStatus: http.StatusOK,
// 		// Middlewares:   middlewares,
// 	}, handler.GetItemPrice)

// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "create-order",
// 		Method:        http.MethodPost,
// 		Path:          paths.Order,
// 		Summary:       "Create Order",
// 		Tags:          []string{tag},
// 		DefaultStatus: http.StatusOK,
// 		// Middlewares:   middlewares,
// 	}, handler.CreateOrder)
// }

// func (h *TecluMobilityHandler) GetItemPrice(ctx context.Context, i *dto.TecluMobilityRequestItemPrice) (
// 	*dto.EntityResponse[entity.ItemPrice], error) {
// 	res, err := h.tecluMobilityService.GetItemPrice(ctx, i)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Fail to get item price", err)
// 	}
// 	var response dto.EntityResponse[entity.ItemPrice]
// 	response.Body.Result = res
// 	return &response, nil
// }

// func (h *TecluMobilityHandler) CreateOrder(ctx context.Context, i *dto.TecluMobilityOrderRequest) (
// 	*dto.TecluMobilityOrderResponse, error,
// ) {
// 	var req common.RequestContext
// 	// companyUuid := "a4b454fa-dc9f-450f-b6fd-d70122052b29"

// 	// company, err := h.companyService.GetCompanyByUuid(ctx, companyUuid)
// 	// if err != nil {
// 	// 	return nil, huma.Error400BadRequest("No company Found", err)
// 	// }
// 	req.Ctx = ctx
// 	// req.ActiveCompany = company
// 	req.LanguageCode = common.LanguageCode(i.AcceptLanguage)
// 	var clientRequest dto.CreateClientRequest
// 	clientRequest.Body.ClientRequestDto = i.Body.ClientRequestDto
// 	deleteAtClient := sql.NullTime{
// 		Valid: true,
// 		Time:  time.Now(),
// 	}
// 	clientRequest.Body.ClientRequestDto.DeleteAt = gorm.DeletedAt(deleteAtClient)
// 	client, err := h.clientService.CreateCustomer(&req, &clientRequest)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest(
// 			h.locale.MustLocalize(
// 				helpers.OptionsLocale.WithID("Error.FailToCreateOrder"),
// 				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
// 			), err)
// 	}

// 	var createAddressPartyRequest dto.CreateAddressPartyRequest
// 	createAddressPartyRequest.Body.Address.FullName = client.GivenName + " " + client.FamilyName
// 	createAddressPartyRequest.Body.Address.CountryCode = i.Body.ClientRequestDto.Country.Code
// 	createAddressPartyRequest.Body.Address.City = i.Body.BillingData.City
// 	createAddressPartyRequest.Body.Address.StreetLine1 = i.Body.BillingData.Address
// 	createAddressPartyRequest.Body.Address.Province = i.Body.BillingData.Estado
// 	createAddressPartyRequest.Body.Address.PostalCode = i.Body.BillingData.PostalCode
// 	createAddressPartyRequest.Body.Address.IdentificationNumber = i.Body.BillingData.TaxIdentificationNumber
// 	createAddressPartyRequest.Body.Address.PhoneNumber = i.Body.ClientRequestDto.PhoneNumber
// 	createAddressPartyRequest.Body.Address.Company = i.Body.ClientRequestDto.CompanyName
// 	createAddressPartyRequest.Body.AddressParty.PartyID = client.ID
// 	createAddressPartyRequest.Body.AddressParty.IsBillingAddress = true
// 	partyAddress, err := h.addressService.CreatePartyAddress(&req, &createAddressPartyRequest)

// 	var createSalesOrderBody dto.CreateSalesOrderBody
// 	salesItemLines, err := h.tecluMobilityService.ToSalesOrderLineDto(ctx, i.Body.OrderData)
// 	createSalesOrderBody.SalesItemLines = salesItemLines
// 	createSalesOrderBody.ClientID = client.ID
// 	createSalesOrderBody.OrderType = entity.ORDER_TYPE_PURCHASE
// 	createSalesOrderBody.DeliveryDate = time.Now()
// 	createSalesOrderBody.Plugins = []string{config.PLUGIN_SQUARE}
// 	createSalesOrderBody.BillingAddressID = &partyAddress.AddressID

// 	orderDeleteAt := sql.NullTime{
// 		Time:  time.Now(),
// 		Valid: true,
// 	}
// 	createSalesOrderBody.DeleteAt = gorm.DeletedAt(orderDeleteAt)
// 	salesOrder, err := h.orderService.CreateSaleOrderService(&req, &createSalesOrderBody)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest(
// 			h.locale.MustLocalize(
// 				helpers.OptionsLocale.WithID("Error.FailToCreateOrder"),
// 				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
// 			), err)
// 	}

// 	var response dto.TecluMobilityOrderResponse
// 	response.Body.PaymentUrl = salesOrder.Body.Result.Data
// 	return &response, nil
// }

// func (h TecluMobilityHandler) SignIn(ctx echo.Context) error {
// 	token ,err := h.jwtService.GenerateToken(common.Claims{
// 		ID: 1,
// 		Uuid: "5f4a8a86-0fd0-4c71-8615-f44a8475d595",
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	time.Sleep(time.Duration(2)* time.Second)
// 	return ctx.JSON(http.StatusOK, SignInResponse{
// 		AccessToken: token,
// 	})
// }
