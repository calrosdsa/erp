package squareservice

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"erp/api/common"
// 	"erp/gen/db/model"
// 	"erp/internal/app/config"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	dtosquare "erp/internal/app/plugin/square/api_square/dto_square"
// 	entitysquare "erp/internal/app/plugin/square/entitiy_square"
// 	squaretypes "erp/internal/app/plugin/square/square_types"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/google/uuid"
// )

// type SquareService struct {
// 	conn          *connection.Connection
// 	configService *config.ConfigService
// 	timeout       *time.Duration
// }

// func NewSquareService(conn *connection.Connection, configService *config.ConfigService,
// 	timeout *time.Duration) *SquareService {
// 	return &SquareService{
// 		conn:          conn,
// 		configService: configService,
// 		timeout:       timeout,
// 	}
// }

// func (s *SquareService) CreateObjectItem(planVariation *squaretypes.CreateSubscriptionPlanResponse,
// 	subsPlan *squaretypes.SubscriptionPlanResponse, item *model.Item) (err error) {
// 	// var objectItem entitysquare.SquareObject
// 	// objectItem.ItemGroupId = item.ItemGroupID
// 	// objectItem.ObjectVariationId = planVariation.CatalogObject.ID
// 	// objectItem.ObjectId = subsPlan.CatalogObject.ID
// 	// objectItem.ItemPriceID = item.ItemPrice.ID
// 	// err = s.conn.Db.Save(&objectItem).Error
// 	return
// }

// func (s *SquareService) CreateSubscriptionSquare(customer *squaretypes.SquareCustomerResponse,
// 	credentails *squaretypes.SquareCredentials, squareMetadata *squaretypes.SquareCustomerMetadata) (err error) {
// 	// orderTempalte,err := s.createSubscriptionOrder(squareMetadata,credentails)
// 	// if err != nil {
// 	// 	return
// 	// }
// 	card, err := s.createCardCustomer(customer, credentails, squareMetadata)
// 	if err != nil {
// 		return
// 	}
// 	res, err := s.createSubscription(customer, &card, credentails, squareMetadata)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(res)

// 	return
// }

// func (s *SquareService) createSubscription(customer *squaretypes.SquareCustomerResponse, card *squaretypes.SquareCardResponse,
// 	credentials *squaretypes.SquareCredentials, metadata *squaretypes.SquareCustomerMetadata) (
// 	res squaretypes.SquareSubscriptionResponse, err error) {
// 	idempotencyKey := uuid.NewString()
// 	body := []byte(fmt.Sprintf(`{
//     "idempotency_key": "%s",
//     "customer_id": "%s",
//     "location_id": "%s",
//     "plan_variation_id": "%s",
//     "card_id": "%s"
//   }`, idempotencyKey, customer.Customer.ID, credentials.LocationId, metadata.PlanVariationId, card.Card.ID))
// 	fmt.Println("CREATE SUBSCRIPTION BODY", string(body))

// 	url := s.GetBaseUrl() + "subscriptions"

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}
// 	// Set headers
// 	req.Header.Set("Square-Version", credentials.ApiVersion)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	req.Header.Set("Content-Type", "application/json")

// 	// Create an HTTP client and send the request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()
// 	// data, err := io.ReadAll(resp.Body)

// 	fmt.Println("CREATE SUBSCRIPTION", resp.Status)
// 	// fmt.Println(string(data), err)
// 	err = json.NewDecoder(resp.Body).Decode(&res)
// 	if err != nil {
// 		return
// 	}
// 	var squareSubscription entitysquare.SquareSubscription
// 	squareSubscription.SubscriptionId = res.Subscription.ID
// 	squareSubscription.CustomerId = res.Subscription.CustomerID
// 	squareSubscription.PlanVariationId = res.Subscription.PlanVariationID
// 	squareSubscription.Status = res.Subscription.Status

// 	fmt.Println(squareSubscription)
// 	err = s.conn.Db.Save(&squareSubscription).Error
// 	if err != nil {
// 		fmt.Println("FAIL TO SAVE SQUARE SUBSCRIPTION")
// 	}
// 	fmt.Println("AFTER", squareSubscription)

// 	return
// }

// func (s *SquareService) createCardCustomer(customer *squaretypes.SquareCustomerResponse, credentials *squaretypes.SquareCredentials,
// 	squareMetadata *squaretypes.SquareCustomerMetadata) (res squaretypes.SquareCardResponse, err error) {
// 	body := []byte(fmt.Sprintf(`{
//     "card": {
//       "customer_id": "%s"
//     },
//     "idempotency_key": "%s",
//     "source_id": "%s"  
//   }
//   `, customer.Customer.ID, squareMetadata.CardRequest.IdEmpotencyKey, squareMetadata.CardRequest.SourceId))

// 	url := s.GetBaseUrl() + "cards"

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}
// 	// Set headers
// 	req.Header.Set("Square-Version", credentials.ApiVersion)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	req.Header.Set("Content-Type", "application/json")

// 	// Create an HTTP client and send the request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()
// 	// data,err := io.ReadAll(resp.Body)
// 	// fmt.Println(string(data),err)
// 	err = json.NewDecoder(resp.Body).Decode(&res)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println("CREATE CARD CUSTOMER", resp.Status, res)
// 	return
// }

// // func (s *SquareService) createSubscriptionOrder(squareMetadata *squaretypes.SquareCustomerMetadata,
// // 	credentials *squaretypes.SquareCredentials) (res squaretypes.SquareOrderResponse,err error) {
// // 	idempotencyKey := uuid.NewString()
// // 	body :=[]byte(fmt.Sprintf(`{
// //     "idempotency_key": "%s",
// //     "order": {
// //       "location_id": "LE40N37TVF5FT",
// //       "state": "DRAFT",
// //       "line_items": [
// //         {
// //           "quantity": "1",
// //           "catalog_object_id": "%s"
// //         }
// //       ],
// //     }
// //   `,idempotencyKey,squareMetadata.ObjectId))
// //   url := s.GetBaseUrl() + "orders"

// // 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// // 	if err != nil {
// // 		fmt.Println("Error creating request:", err)
// // 		return
// // 	}
// // 	// Set headers
// // 	req.Header.Set("Square-Version", credentials.ApiVersion)
// // 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// // 	req.Header.Set("Content-Type", "application/json")

// // 	// Create an HTTP client and send the request
// // 	client := &http.Client{}
// // 	resp, err := client.Do(req)
// // 	if err != nil {
// // 		return
// // 	}
// // 	defer resp.Body.Close()
// // 	fmt.Println(resp.Status)
// // 	// data,err := io.ReadAll(resp.Body)
// // 	// fmt.Println(string(data),err)
// // 	err = json.NewDecoder(resp.Body).Decode(&res)
// // 	if err != nil {
// // 		return
// // 	}
// // 	return
// // }

// func (s *SquareService) RetrieveObjectRequest(ctx context.Context, d *dtosquare.SquareObjectRequest) (
// 	res dtosquare.SquareObjectResponse, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, *s.timeout)
// 	defer cancel()
// 	var itemGroup entity.ItemGroup
// 	if err = s.conn.Db.WithContext(ctx).First(&itemGroup, "uuid = $1", d.ItemGroupUuid).Error; err != nil {
// 		return
// 	}
// 	credentials, err := s.getCredentials(ctx, itemGroup.CompanyID)
// 	if err != nil {
// 		return
// 	}
// 	// var squareObjectItem entitysquare.ObjectItemSquare
// 	// if err = s.conn.Db.WithContext(ctx).First(&squareObjectItem, "item_group_id = $1", itemGroup.ID).Error; err != nil {
// 	// 	return
// 	// }
// 	// fmt.Println(squareObjectItem.ObjectId)
// 	url := fmt.Sprintf("%scatalog/object/%s", s.GetBaseUrl(), d.ObjectId)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}

// 	// Set headers
// 	req.Header.Set("Square-Version", credentials.ApiVersion)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	req.Header.Set("Content-Type", "application/json")

// 	// Create an HTTP client and send the request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.Status)
// 	// data,err := io.ReadAll(resp.Body)
// 	// fmt.Println(string(data),err)
// 	err = json.NewDecoder(resp.Body).Decode(&res.Body.PlanVariation)
// 	if err != nil {
// 		return
// 	}

// 	err = s.conn.Db.WithContext(ctx).Where("object_variation_id = ?", d.ObjectId).
// 		Preload("ItemPrice.Tax").Preload("ItemPrice.ItemPriceList").First(&res.Body.SquareObject).Error
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// func (s *SquareService) RetrieveCatalagoRequest(ctx context.Context, d *dtosquare.SquareCatalogRequest) (
// 	res dtosquare.SquareCatalogResponse, err error,
// ) {
// 	ctx, cancel := context.WithTimeout(ctx, *s.timeout)
// 	defer cancel()
// 	var itemGroup entity.ItemGroup
// 	if err = s.conn.Db.WithContext(ctx).First(&itemGroup, "uuid = $1", d.ItemGroupUuid).Error; err != nil {
// 		return
// 	}
// 	credentials, err := s.getCredentials(ctx, itemGroup.CompanyID)
// 	if err != nil {
// 		return
// 	}
// 	var squareObjectItem entitysquare.SquareObject
// 	if err = s.conn.Db.WithContext(ctx).First(&squareObjectItem, "item_group_id = $1", itemGroup.ID).Error; err != nil {
// 		return
// 	}
// 	fmt.Println(squareObjectItem.ObjectId)
// 	url := fmt.Sprintf("%scatalog/object/%s", s.GetBaseUrl(), squareObjectItem.ObjectId)

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}

// 	// Set headers
// 	req.Header.Set("Square-Version", credentials.ApiVersion)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	req.Header.Set("Content-Type", "application/json")

// 	// Create an HTTP client and send the request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.Status)
// 	// data,err := io.ReadAll(resp.Body)
// 	// fmt.Println(string(data),err)
// 	err = json.NewDecoder(resp.Body).Decode(&res.Body.Catalog)
// 	if err != nil {
// 		return
// 	}

// 	// var squareObjects []entitysquare.SquareObject
// 	err = s.conn.Db.WithContext(ctx).Where("item_group_id = ?", itemGroup.ID).
// 		Preload("ItemPrice.Tax").Find(&res.Body.Objects).Error
// 	if err != nil {
// 		return
// 	}

// 	return
// }

// func (s *SquareService) getCredentials(ctx context.Context, companyId int64) (squareCredential squaretypes.SquareCredentials, err error) {
// 	var companyPlugin entity.CompanyPlugins
// 	pass := s.configService.GetDbConfig().CryptoPass
// 	err = s.conn.Db.WithContext(ctx).Raw("select company_id,plugin, pgp_sym_decrypt(credentials::bytea, $1) as credentials from company_plugins where company_id = $2 and plugin = $3",
// 		pass, companyId, config.PLUGIN_SQUARE).
// 		Scan(&companyPlugin).Error
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(companyPlugin.Credentials)
// 	err = json.Unmarshal([]byte(companyPlugin.Credentials), &squareCredential)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// func (s *SquareService) GetCredentials(req *common.RequestContext) (squareCredential squaretypes.SquareCredentials, err error) {
// 	var companyPlugin entity.CompanyPlugins
// 	pass := s.configService.GetDbConfig().CryptoPass
// 	err = s.conn.Db.Raw("select company_id,plugin, pgp_sym_decrypt(credentials::bytea, $1) as credentials from company_plugins where company_id = $2 and plugin = $3",
// 		pass, req.ActiveCompany.ID, config.PLUGIN_SQUARE).
// 		Scan(&companyPlugin).Error
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(companyPlugin.Credentials)
// 	err = json.Unmarshal([]byte(companyPlugin.Credentials), &squareCredential)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// func (s *SquareService) GetBaseUrl() string {
// 	return "https://connect.squareupsandbox.com/v2/"
// }
