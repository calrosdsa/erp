package eventsquarehandler

// import (
// 	"context"
// 	"erp/api/common"
// 	"erp/api/dto"
// 	"erp/gen/db/model"
// 	"erp/internal/app/config"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	"erp/internal/app/event-bus/event"
// 	entitysquare "erp/internal/app/plugin/square/entitiy_square"
// 	"erp/internal/app/service/helpers"
// 	"erp/internal/app/service/services"
// 	clientservice "erp/internal/app/service/services/client_service"
// 	"erp/internal/app/service/services/company_service"
// 	sellingservice "erp/internal/app/service/services/selling_service"
// 	"erp/pkg/logger"
// 	"fmt"
// 	"time"

// 	"github.com/asaskevich/EventBus"
// )

// type EventSquareHandler struct {
// 	conn              *connection.Connection
// 	salesOrderService *sellingservice.SalesOrderService
// 	companyService    *company_service.CompanyService
// 	timeout           *time.Duration
// 	currencyHelper    *helpers.CurrencyHelper
// 	emitLog           helpers.EmitLog
// 	configService     *config.ConfigService
// 	clientService     *clientservice.ClientService
// 	bus               EventBus.Bus
// }

// func NewEventSquareHandler(conn *connection.Connection, busP *EventBus.Bus, services *services.Services,
// 	configService *config.ConfigService, helpers *helpers.Helpers) {
// 	bus := *busP
// 	timeout := configService.GetTimeoutAPICall()
// 	handler := &EventSquareHandler{
// 		salesOrderService: services.SalesOrderService,
// 		conn:              conn,
// 		timeout:           &timeout,
// 		companyService:    services.CompanyService,
// 		currencyHelper:    helpers.Currency,
// 		emitLog:           helpers.Logger.EmitLog("eventsquare-handler"),
// 		configService:     configService,
// 		bus:               bus,
// 		clientService:     services.ClientService,
// 	}
// 	bus.Subscribe(event.SQUARE_PAYMENT_COMPLETED_EVENT, handler.SquarePaymentCompleted)
// }

// func (e *EventSquareHandler) SquarePaymentCompleted(d event.SquarePaymenrCompleted) {
// 	ctx, cancel := context.WithTimeout(context.Background(), *e.timeout)
// 	defer cancel()

// 	//IF EXIST ORDER UPDATED OTHERWISE CREATE ONE BASE ON SUBSCRIPTION
// 	orderExist := e.validateIfSalesOrderExist(ctx, d.Body.Data.Object.Payment.OrderID)
// 	if orderExist {
// 		return
// 	}
// 	createdOrderRequest := dto.CreateSalesOrderBody{}
// 	squareCustomer, err := e.getSquareCustomer(ctx, d.Body.Data.Object.Payment.CustomerID)
// 	if err != nil {
// 		e.emitLog.Err(err, logger.OptionsLog.WithOperation("getSquareCustomer"))
// 		// fmt.Println("fail to get customer", err)
// 		return
// 	}
// 	squareSubscription, err := e.getLastActiverSubscription(ctx, squareCustomer.CustomerId)
// 	if err != nil {
// 		e.emitLog.Err(err, logger.OptionsLog.WithOperation("getLastActiverSubscription"))
// 		// fmt.Println("Fail to get square subscription", err)
// 		return
// 	}
// 	squareObject, err := e.getSquareObject(ctx, squareSubscription.PlanVariationId)
// 	if err != nil {
// 		e.emitLog.Err(err, logger.OptionsLog.WithOperation("getSquareObject"))
// 		// fmt.Println("FAIL TO GET SQUARE OBJECT", err)
// 		return
// 	}
// 	// company, err := e.companyService.GetCompanyByUuid(ctx, "")
// 	// if err != nil {
// 	// 	e.emitLog.Err(err, logger.OptionsLog.WithOperation("GetCompany"))
// 	// 	// fmt.Println("FAIL TO GET COMPANY", err)
// 	// 	return
// 	// }
// 	var req common.RequestContext
// 	// req.ActiveCompany = company
// 	createdOrderRequest.ClientID = squareCustomer.PartyID
// 	createdOrderRequest.DeliveryDate = time.Now()
// 	createdOrderRequest.OrderType = entity.ORDER_TYPE_SERVICE

// 	var salesItemLine dto.SalesItemLineDto
// 	salesItemLine.Currency = squareObject.ItemPrice.ItemPriceList.Currency
// 	salesItemLine.ItemPriceID = squareObject.ItemPriceID
// 	salesItemLine.ItemQuanitity = squareObject.ItemPrice.ItemQuantity
// 	salesItemLine.Rate = e.currencyHelper.IntToFloat(squareObject.ItemPrice.Rate)
// 	createdOrderRequest.SalesItemLines = []dto.SalesItemLineDto{salesItemLine}
// 	createdOrderRequest.Plugins = []string{config.PLUGIN_SQUARE}
// 	req.Ctx = ctx
// 	_, err = e.salesOrderService.CreateSaleOrderService(&req, &createdOrderRequest)
// 	if err != nil {
// 		e.emitLog.Err(err, logger.OptionsLog.WithOperation("CreateSaleOrderService"))
// 		// fmt.Println("Fail to create order", err)
// 	}
// }

// func (e *EventSquareHandler) validateIfSalesOrderExist(ctx context.Context, orderId string) bool {
// 	var squareOrder entitysquare.SquareOrder
// 	fmt.Println("BEFORE PRELOAD")
// 	err := e.conn.Db.WithContext(ctx).Where(&entitysquare.SquareOrder{SquareOrderId: orderId}).
// 		First(&squareOrder).Error
// 	if err != nil {
// 		return false
// 	}
// 	err = e.conn.Db.WithContext(ctx).Exec("UPDATE clients SET deleted_at = NULL WHERE id = ?",
// 		squareOrder.PartyID).Error
// 	if err != nil {
// 		return false
// 	}
// 	err = e.conn.Db.WithContext(ctx).Exec("UPDATE sales_orders SET deleted_at = NULL WHERE id = ?",
// 		squareOrder.SalesOrderID).Error
// 	if err != nil {
// 		return false
// 	}
// 	var client entity.Client
// 	fmt.Println("CLIENTID", squareOrder.PartyID)
// 	err = e.conn.Db.WithContext(ctx).Where(&entity.Client{ID: squareOrder.PartyID}).
// 		Preload("Company").First(&client).Error

// 	fmt.Println("AFTER PRELOAD", squareOrder)

// 	//SENDING CLIENT EMAIL INVITATION NOTIFICATION
// 	fmt.Println("SENDING N FIRST", squareOrder)
// 	reqCtx := e.createRequestContext(&client.Company)
// 	fmt.Println("SENDING N FIRST 2", squareOrder)
// 	e.clientService.SendCredentialEmail(reqCtx, client)

// 	fmt.Println("SENDING N AFTER", squareOrder)

// 	defer func() {
// 		if err != nil {
// 			e.emitLog.Err(err, logger.OptionsLog.WithMethod("validateIfSalesOrderExist"))
// 		}
// 	}()
// 	return true
// }

// func (e *EventSquareHandler) createRequestContext(company *model.Company) common.RequestContext {
// 	var reqCtx common.RequestContext
// 	reqCtx.LanguageCode = common.LanguageCode(e.configService.GetDefaultLanguage())
// 	// reqCtx.ActiveCompany = *company
// 	return reqCtx
// }

// func (e *EventSquareHandler) getSquareCustomer(ctx context.Context, customerId string) (res entitysquare.SquareCustomer, err error) {
// 	err = e.conn.Db.WithContext(ctx).Where("customer_id = ?", customerId).Preload("Client.Company").First(&res).Error
// 	return
// }

// func (e *EventSquareHandler) getSquareObject(ctx context.Context, planVariationId string) (res entitysquare.SquareObject, err error) {
// 	err = e.conn.Db.WithContext(ctx).Where("object_variation_id = ?", planVariationId).
// 		Preload("ItemPrice.ItemPriceList").First(&res).Error
// 	return
// }

// func (e *EventSquareHandler) getLastActiverSubscription(ctx context.Context, customerId string) (res entitysquare.SquareSubscription, err error) {
// 	err = e.conn.Db.WithContext(ctx).Where("customer_id = ? and status = ?", customerId, entitysquare.SQUARE_ACTIVE_SUBSCRIPTION).
// 		Order("created_at DESC").
// 		First(&res).Error
// 	return
// }
