package squarestrategies

// import (
// 	"bytes"
// 	"encoding/json"
// 	"erp/api/common"
// 	"erp/gen/db/model"
// 	"erp/internal/app/config"
// 	_strategy "erp/internal/app/config/stock_config"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	squareservice "erp/internal/app/plugin/square/square_service"
// 	squaretypes "erp/internal/app/plugin/square/square_types"
// 	"erp/internal/app/service/helpers"
// 	"erp/pkg/logger"
// 	"fmt"
// 	"net/http"
// )

// var ( //interface comformance checks
// 	_ _strategy.ItemStrategy = (*SquareItemStrategy)(nil)
// )

// // const SQUARE_URL = "https://connect.squareupsandbox.com/v2/"

// type SquareItemStrategy struct {
// 	conn          *connection.Connection
// 	configService *config.ConfigService
// 	squareService *squareservice.SquareService
// 	taxHelper     *helpers.TaxHelper
// 	emitLog helpers.EmitLog
// }

// func NewSquareItemStrategy(conn *connection.Connection, configService *config.ConfigService,
// 	squareService *squareservice.SquareService, helpers *helpers.Helpers) _strategy.ItemStrategy {
// 	// credentials :=
// 	return &SquareItemStrategy{
// 		conn:          conn,
// 		configService: configService,
// 		squareService: squareService,
// 		taxHelper: helpers.Tax,
// 		emitLog: helpers.Logger.EmitLog("square-item-strategy"),
// 	}
// }

// func (s *SquareItemStrategy) CreateItem(req *common.RequestContext, item *model.Item, itemPrice *entity.ItemPrice) error {
// 	credentials, err := s.squareService.GetCredentials(req)
// 	if err != nil {
// 		s.emitLog.Err(err,logger.OptionsLog.WithMethod("CreateItem_GetCredentials"))
// 		// fmt.Println("Fail to get credentials", err)
// 	}
// 	// fmt.Println("CREDENTIALS", credentials)
// 	subsPlan, err := s.createSubscriptionPlan(credentials, item)
// 	if err != nil {
// 		s.emitLog.Err(err,logger.OptionsLog.WithMethod("CreateItem_createSubscriptionPlan"))
// 		return err
// 	}
// 	planVariation, err := s.createSubscriptionPlanVariation(subsPlan, credentials, item, itemPrice)
// 	if err != nil {
// 		s.emitLog.Err(err,logger.OptionsLog.WithMethod("CreateItem_createSubscriptionPlanVariation"))
// 		return err
// 	}
// 	err = s.squareService.CreateObjectItem(&planVariation, &subsPlan, item)
// 	if err != nil {
// 		s.emitLog.Err(err,logger.OptionsLog.WithMethod("CreateItem_CreateObjectItem"))
// 		// fmt.Println("Fail to save object item", err)
// 		return err
// 	}
// 	fmt.Println("CATALOG OBJECT ID", subsPlan.CatalogObject.ID)
// 	fmt.Println("PLAN VARIATION CATALOG OBJECT ID", planVariation.CatalogObject.ID)
// 	return err
// }

// // func (s *SquareItemStrategy) batchUpsert(){}

// func (s *SquareItemStrategy) createSubscriptionPlan(d squaretypes.SquareCredentials, item *model.Item) (squaretypes.SubscriptionPlanResponse, error) {
// 	url := s.squareService.GetBaseUrl() + "catalog/object"
// 	body := []byte(fmt.Sprintf(`{
//     "idempotency_key": "%s",
//     "object": {
//       "type": "SUBSCRIPTION_PLAN",
//       "id": "#%d",
//       "subscription_plan_data": {
//         "name": "%s",
//         "all_items": true
//       }
//     }
//   }`, item.ItemGroup.Uuid, item.ItemGroup.ID, item.ItemGroup.Name))
// 	fmt.Println(string(body))
// 	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return squaretypes.SubscriptionPlanResponse{}, err
// 	}
// 	r.Header.Add("Content-Type", "application/json")
// 	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", d.AccessToken))
// 	r.Header.Add("Square-Version", d.ApiVersion)
// 	client := &http.Client{}
// 	res, err := client.Do(r)
// 	if err != nil {
// 		return squaretypes.SubscriptionPlanResponse{}, err
// 	}

// 	defer res.Body.Close()
// 	subscriptionPlanResponse := squaretypes.SubscriptionPlanResponse{}
// 	fmt.Println(res.Status)
// 	derr := json.NewDecoder(res.Body).Decode(&subscriptionPlanResponse)
// 	if derr != nil {
// 		return squaretypes.SubscriptionPlanResponse{}, err
// 	}
// 	return subscriptionPlanResponse, err
// }

// func (s *SquareItemStrategy) createSubscriptionPlanVariation(
// 	subsP squaretypes.SubscriptionPlanResponse,
// 	d squaretypes.SquareCredentials,
// 	item *model.Item,
// 	itemPrice *model.ItemPrice,
// ) (
// 	planV squaretypes.CreateSubscriptionPlanResponse, err error) {

// 	url := s.squareService.GetBaseUrl() + "catalog/object"

// 	cadence := s.getCadence(item.UnitOfMeasure.Code)
// 	totalPriceWithTax := s.taxHelper.CalculateTotalWithTax(item.ItemPrice.Rate,item.ItemPrice.Tax.Value)
// 	body := []byte(fmt.Sprintf(`{
//     "idempotency_key": "%s",
//     "object": {
//       "type": "SUBSCRIPTION_PLAN_VARIATION",
//       "id": "#1",
//       "subscription_plan_variation_data": {
//         "name": "%s",
//         "phases": [
//           {
//             "cadence": "%s",
//             "ordinal": 0,
//             "periods": %d,
//             "pricing": {
//               "type": "STATIC",
// 				"price_money": {
//                 "amount": %d,
//                 "currency": "%s"
//               }
//             }
//           }
//         ],
//         "subscription_plan_id": "%s"
//       }
//     }
//   }`, item.Uuid, item.Name, cadence, itemPrice.ItemQuantity,totalPriceWithTax, itemPrice.ItemPriceList.Currency, subsP.CatalogObject.ID))
// 	fmt.Println("Create plan variation", string(body))
// 	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return squaretypes.CreateSubscriptionPlanResponse{}, err
// 	}
// 	r.Header.Add("Content-Type", "application/json")
// 	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", d.AccessToken))
// 	r.Header.Add("Square-Version", d.ApiVersion)
// 	client := &http.Client{}
// 	res, err := client.Do(r)
// 	if err != nil {
// 		return squaretypes.CreateSubscriptionPlanResponse{}, err
// 	}

// 	defer res.Body.Close()
// 	subscriptionPlanVariationResponse := squaretypes.CreateSubscriptionPlanResponse{}
// 	fmt.Println(res.Status)
// 	derr := json.NewDecoder(res.Body).Decode(&subscriptionPlanVariationResponse)
// 	if derr != nil {
// 		return squaretypes.CreateSubscriptionPlanResponse{}, err
// 	}
// 	return subscriptionPlanVariationResponse, err
// }

// // func (s *SquareItemStrategy) getPlanVariation()

// func (s *SquareItemStrategy) getCadence(uomCode string) string {
// 	switch uomCode {
// 	case "MON":
// 		return "MONTHLY"
// 	case "DAY":
// 		return "DAILY"
// 	case "WEE":
// 		return "WEEKLY"
// 	default:
// 		return ""
// 	}
// }
