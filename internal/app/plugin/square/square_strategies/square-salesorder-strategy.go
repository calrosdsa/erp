package squarestrategies

// import (
// 	"bytes"
// 	"encoding/json"
// 	"erp/api/common"
// 	_strategy "erp/internal/app/config/selling_config"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	entitysquare "erp/internal/app/plugin/square/entitiy_square"
// 	squareservice "erp/internal/app/plugin/square/square_service"
// 	squaretypes "erp/internal/app/plugin/square/square_types"
// 	"erp/internal/app/service/helpers"
// 	"erp/pkg/logger"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/google/uuid"
// )

// var (
// 	_ _strategy.SalesOrderStrategy = (*SquareSalesOrderStrategy)(nil)
// )

// type SquareSalesOrderStrategy struct {
// 	conn          *connection.Connection
// 	squareService *squareservice.SquareService
// 	timeout       time.Duration
// 	emitLog       helpers.EmitLog
// }

// func NewSquareSalesOrderStrategy(conn *connection.Connection, squareService *squareservice.SquareService,
// 	timeout time.Duration,
// 	helpers *helpers.Helpers,
// ) _strategy.SalesOrderStrategy {
// 	return &SquareSalesOrderStrategy{
// 		conn:          conn,
// 		squareService: squareService,
// 		timeout:       timeout,
// 		emitLog:       helpers.Logger.EmitLog("square-sales-order-strategy"),
// 	}
// }

// func (s *SquareSalesOrderStrategy) CreateSalesOrder(req *common.RequestContext, salesOrder *entity.SalesOrder, lines []entity.SalesItemLine) (
// 	res string, err error) {
// 	if len(lines) == 0 {
// 		err = fmt.Errorf("No line items in order")
// 		return
// 	}
// 	credentials, err := s.squareService.GetCredentials(req)
// 	if err != nil {
// 		return
// 	}
// 	squareOrderResp, err := s.createPaymentLink(&credentials, lines)
// 	if err != nil {
// 		fmt.Println("Fail to create square order", err)
// 	}
// 	err = s.createSquareOrder(&squareOrderResp, salesOrder)
// 	if err != nil {
// 		return
// 	}

// 	return squareOrderResp.PaymentLink.URL, err
// }

// func (s *SquareSalesOrderStrategy) createSquareOrder(orderResp *squaretypes.SquareOrderResponse, salesOrder *entity.SalesOrder) (
// 	err error) {
// 	var squareOrder entitysquare.SquareOrder
// 	squareOrder.SalesOrderID = salesOrder.ID
// 	squareOrder.SquareOrderId = orderResp.PaymentLink.OrderID
// 	squareOrder.PartyID = salesOrder.PartyID
// 	err = s.conn.Db.Save(&squareOrder).Error
// 	return
// }

// func (s *SquareSalesOrderStrategy) getTotalOrderFromLines(lines []entity.SalesItemLine) ([]squaretypes.LineItem, error) {
// 	lineItems := make([]squaretypes.LineItem, len(lines))
// 	for i, line := range lines {
// 		lineItem := squaretypes.LineItem{}
// 		var item model.Item
// 		err := s.conn.Db.Where(&model.Item{Base: entity.Base{ID: line.ItemPrice.ItemID}}).First(&item).Error
// 		if err != nil {
// 			s.emitLog.Err(err, logger.OptionsLog.WithMethod("getTotalOrderFromLines"))
// 			return lineItems, err
// 		}
// 		lineItem.Name = item.Name
// 		lineItem.Quantity = strconv.Itoa(int(line.ItemQuantity))
// 		lineItem.BasePriceMoney.Amount = line.Rate
// 		lineItem.BasePriceMoney.Currency = line.Currency
// 		lineItems[i] = lineItem
// 	}
// 	// for _,line := range lines{
// 	// 	res += line.Rate
// 	// }
// 	return lineItems, nil
// }

// func (s *SquareSalesOrderStrategy) createPaymentLink(credentials *squaretypes.SquareCredentials, lines []entity.SalesItemLine) (
// 	res squaretypes.SquareOrderResponse, err error) {
// 	url := s.squareService.GetBaseUrl() + "online-checkout/payment-links"
// 	client := http.Client{
// 		Timeout: s.timeout,
// 	}
// 	idempotencyKey := uuid.NewString()
// 	lineItems, err := s.getTotalOrderFromLines(lines)
// 	if err != nil {
// 		return
// 	}
// 	lineItemsStr, err := json.Marshal(lineItems)
// 	if err != nil {
// 		return
// 	}

// 	body := []byte(fmt.Sprintf(`{
//     "idempotency_key": "%s",
//     "order" : {
//       "location_id": "%s",
// 	  "line_items":%s
//     }
//   }`, idempotencyKey, credentials.LocationId, string(lineItemsStr)))
// 	fmt.Println("BODY", string(body))
// 	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return
// 	}
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	req.Header.Set("Square-Version", credentials.ApiVersion)
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("FAILT TO EXECUTE REQUEST", err)
// 		return
// 	}
// 	fmt.Println("create payment link", resp.StatusCode)
// 	if resp.StatusCode != http.StatusOK {
// 		err = fmt.Errorf("Fail to create payment link")
// 		return
// 	}
// 	err = json.NewDecoder(resp.Body).Decode(&res)
// 	return
// }

// func (s *SquareSalesOrderStrategy) GetSalesOrderDetail(req *common.RequestContext, salesOrder *entity.SalesOrder, lines []entity.SalesItemLine) (
// 	res string, err error) {
// 	if salesOrder.OrderType == entity.ORDER_TYPE_SERVICE {
// 		res, err = s.getSubscriptionData(req, salesOrder, lines)
// 		return
// 	}
// 	return
// }

// func (s *SquareSalesOrderStrategy) getSubscriptionData(req *common.RequestContext, salesOrder *entity.SalesOrder, lines []entity.SalesItemLine) (
// 	res string, err error) {
// 	squareObject, err := s.getSquareObject(req, lines)
// 	if err != nil {
// 		return
// 	}
// 	squareCustomer, err := s.getSquareCustomer(req, salesOrder)
// 	if err != nil {
// 		return
// 	}
// 	squareSubscription, err := s.getSquareSubscription(req, squareCustomer.CustomerId, squareObject.ObjectVariationId)
// 	if err != nil {
// 		return
// 	}
// 	subscription, err := s.getSubscription(req, &squareSubscription)
// 	data := squaretypes.SalesOrderSquareSubscription{}
// 	data.Subscription = subscription
// 	data.SquareSubscription = squareSubscription
// 	byteData, err := json.Marshal(data)
// 	return string(byteData), err
// }

// func (s *SquareSalesOrderStrategy) getSubscription(r *common.RequestContext, squareSubscription *entitysquare.SquareSubscription) (
// 	res squaretypes.SquareSubscriptionResponse, err error) {
// 	url := s.squareService.GetBaseUrl() + fmt.Sprintf("subscriptions/%s?include=actions", squareSubscription.SubscriptionId)
// 	credentials, err := s.squareService.GetCredentials(r)
// 	if err != nil {
// 		return
// 	}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request", err)
// 		return
// 	}

// 	// Set the required headers
// 	req.Header.Set("Square-Version", credentials.ApiVersion)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	req.Header.Set("Content-Type", "application/json")

// 	// Send the request using the default client
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error making request", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	res.Subscription.Actions = make([]squaretypes.SubscriptionActions, 0)
// 	// Read and print the response body
// 	derr := json.NewDecoder(resp.Body).Decode(&res)
// 	if derr != nil {
// 		return
// 	}
// 	return
// }

// func (s *SquareSalesOrderStrategy) getSquareSubscription(req *common.RequestContext, customerId string, planVariationID string) (
// 	res entitysquare.SquareSubscription, err error) {
// 	s.conn.Db.WithContext(req.Ctx).First(&res, "customer_id = ? and plan_variation_id = ?", customerId, planVariationID)
// 	return
// }

// func (s *SquareSalesOrderStrategy) getSquareCustomer(req *common.RequestContext, salesOrder *entity.SalesOrder) (
// 	res entitysquare.SquareCustomer, err error) {
// 	// var res entitysquare.SquareObject
// 	s.conn.Db.WithContext(req.Ctx).First(&res, "client_id = ?", salesOrder.PartyID)
// 	return
// }

// func (s *SquareSalesOrderStrategy) getSquareObject(req *common.RequestContext, lines []entity.SalesItemLine) (
// 	squareObject entitysquare.SquareObject, err error) {
// 	if len(lines) == 0 {
// 		err = fmt.Errorf("NO line items present")
// 		return
// 	}
// 	// var squareObject entitysquare.SquareObject
// 	s.conn.Db.WithContext(req.Ctx).First(&squareObject, "item_price_id = ?", lines[0].ItemPriceID)
// 	return
// }
