package handlersquare

// import (
// 	"context"
// 	"erp/api/common"
// 	"erp/api/dto"
// 	"erp/gen/db/model"
// 	"erp/internal/app/entity"
// 	"erp/internal/app/event-bus/event"
// 	dtosquare "erp/internal/app/plugin/square/api_square/dto_square"
// 	squareservice "erp/internal/app/plugin/square/square_service"
// 	"erp/internal/app/service/helpers"
// 	"fmt"
// 	"net/http"

// 	"github.com/danielgtaylor/huma/v2"
// )

// // optional code omitted

// type SquareHandler struct {
// 	squareService *squareservice.SquareService
// 	eventHelper   *helpers.EventHelper
// }

// func NewHandler(
// 	api *huma.API,
// 	squareService *squareservice.SquareService,
// 	base string,
// 	tag string,
// 	helpers *helpers.Helpers,
// ) {
// 	paths := NewSquarePath(base)
// 	handler := SquareHandler{
// 		squareService: squareService,
// 		eventHelper:   helpers.Event,
// 	}

// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "test-event",
// 		Method:        http.MethodGet,
// 		Path:          paths.Test,
// 		Summary:       "Test Event",
// 		Tags:          []string{tag},
// 		DefaultStatus: http.StatusOK,
// 	}, handler.TestEvent)

// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "get-catalog-square",
// 		Method:        http.MethodGet,
// 		Path:          paths.Catalog,
// 		Summary:       "Get square catalog",
// 		Tags:          []string{tag},
// 		DefaultStatus: http.StatusOK,
// 	}, handler.GetSquareCatalog)

// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "get-object-square",
// 		Method:        http.MethodGet,
// 		Path:          paths.Object,
// 		Summary:       "Get square object",
// 		Tags:          []string{tag},
// 		DefaultStatus: http.StatusOK,
// 	}, handler.GetSquareObject)

// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "payment-webhook",
// 		Method:        http.MethodPost,
// 		Path:          paths.PaymentWeebhook,
// 		Summary:       "Payment square wenhook",
// 		Tags:          []string{tag},
// 		DefaultStatus: http.StatusOK,
// 	}, handler.PaymentWeebhook)
// }

// func (h *SquareHandler) TestEvent(ctx context.Context, i *struct{}) (*struct{}, error) {
// 	h.eventHelper.Publish(event.NOTIFICATION_EVENT, event.NotificationData{
// 		NotificationEventType: event.NOTIFY_NEW_CLIENT,
// 		Data:event.NotificationPayload{
// 			Payload: entity.Client{
// 				GivenName:    "Jorge",
// 				FamilyName:   "Miranda",
// 				EmailAddress: "jorgemiranda0180@gmail.com",
// 				Company: model.Company{
// 					Name: "CUATROPF",
// 				},
// 				UserID: 16,
// 			},
// 			RequestContext: common.RequestContext{
// 				LanguageCode: common.LanguageCodeEN,
// 				ActiveCompany:  model.Company{
// 					Name: "CUATROPF",
// 				},
// 			},
// 		},
// 	})
// 	return nil, nil
// }

// func (h *SquareHandler) PaymentWeebhook(ctx context.Context, i *dtosquare.PaymentWeebhookRequest) (*dto.ResponseMessage, error) {
// 	var response dto.ResponseMessage
// 	response.Body.Message = "Success"
// 	fmt.Println("PAYMENT RESPONSE", i)
// 	if i.Body.PaymentBody.Data.Object.Payment.Status == "COMPLETED" {
// 		h.eventHelper.Publish(event.SQUARE_PAYMENT_COMPLETED_EVENT, event.SquarePaymenrCompleted{
// 			Body: i.Body.PaymentBody,
// 		})
// 	}

// 	return &response, nil
// }

// func (h *SquareHandler) GetSquareCatalog(ctx context.Context, i *dtosquare.SquareCatalogRequest) (*dtosquare.SquareCatalogResponse, error) {
// 	// response := dtosquare.SquareCatalogResponse{}

// 	res, err := h.squareService.RetrieveCatalagoRequest(ctx, i)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Fail to fetch square catalog", err)
// 	}
// 	// fmt.Println("RES", res)
// 	// response = res
// 	return &res, nil
// }

// func (h *SquareHandler) GetSquareObject(ctx context.Context, i *dtosquare.SquareObjectRequest) (*dtosquare.SquareObjectResponse, error) {
// 	// response := dtosquare.SquareObjectResponse{}

// 	res, err := h.squareService.RetrieveObjectRequest(ctx, i)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Fail to fetch square catalog", err)
// 	}
// 	// fmt.Println("RES", res)
// 	// response.Body.PlanVariation = res
// 	return &res, nil
// }
