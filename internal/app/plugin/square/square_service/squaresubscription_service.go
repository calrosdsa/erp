package squareservice

// import (
// 	"context"
// 	"encoding/json"
// 	"erp/api/common"
// 	"erp/internal/app/config"
// 	"erp/internal/app/connection"
// 	dtosquare "erp/internal/app/plugin/square/api_square/dto_square"
// 	entitysquare "erp/internal/app/plugin/square/entitiy_square"
// 	squaretypes "erp/internal/app/plugin/square/square_types"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// type SquareSubscriptionService struct {
// 	conn          *connection.Connection
// 	configService *config.ConfigService
// 	timeout       *time.Duration
// 	squareService *SquareService
// }

// func NewSquareSubscriptionService(
// 	conn *connection.Connection,
// 	configService *config.ConfigService,
// 	timeout *time.Duration,
// 	squareService *SquareService,
// ) *SquareSubscriptionService {
// 	return &SquareSubscriptionService{
// 		conn:          conn,
// 		configService: configService,
// 		timeout:       timeout,
// 		squareService: squareService,
// 	}
// }

// func (s *SquareSubscriptionService) CancelSubscription(req *common.RequestContext, d *dtosquare.RequestSubscriptionCancel) (
// 	res *squaretypes.SquareErrors,err error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
// 	defer cancel()
// 	clientId := req.GetClientID()
// 	squareCustomer, err := s.getSquareCustomer(ctx, clientId)
// 	if err != nil {
// 		return
// 	}

// 	squareSubscription, err := s.getSquareSubscription(ctx, squareCustomer.CustomerId, d.Body.SubscriptionId)
// 	if err != nil {
// 		return
// 	}
// 	credentials, err := s.squareService.GetCredentials(req)
// 	if err != nil {
// 		return
// 	}
// 	url := s.squareService.GetBaseUrl() + fmt.Sprintf("subscriptions/%s/cancel", squareSubscription.SubscriptionId)
// 	request, err := http.NewRequest("POST", url, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}

// 	// Set headers
// 	request.Header.Set("Square-Version", credentials.ApiVersion)
// 	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	request.Header.Set("Content-Type", "application/json")

// 	// Create an HTTP client and send the request
// 	client := &http.Client{
// 		Timeout: *s.timeout,
// 	}
// 	resp, err := client.Do(request)
// 	if err != nil {
// 		fmt.Println("RESP ERROR",err)
// 		return
// 	}
// 	fmt.Println("STATUS",resp.Status)
// 	defer resp.Body.Close()
// 	if resp.StatusCode == http.StatusBadRequest {
// 		err = json.NewDecoder(resp.Body).Decode(&res)
// 		fmt.Println(res)
// 		return
// 	}
// 	if resp.StatusCode >= http.StatusBadRequest {
// 		err = fmt.Errorf("Fail to cancel subscription")
// 		return
// 	}
// 	return
// }

// func (s *SquareSubscriptionService) getSquareSubscription(ctx context.Context, customerId string, subscriptionId string) (
// 	res entitysquare.SquareSubscription, err error) {
// 	err = s.conn.Db.WithContext(ctx).Where(&entitysquare.SquareSubscription{CustomerId: customerId, SubscriptionId: subscriptionId}).
// 		First(&res).Error
// 	return
// }

// func (s *SquareSubscriptionService) getSquareCustomer(ctx context.Context, clientId uint) (res entitysquare.SquareCustomer, err error) {
// 	err = s.conn.Db.WithContext(ctx).First(&res, "client_id = ?", clientId).Error
// 	return
// }
