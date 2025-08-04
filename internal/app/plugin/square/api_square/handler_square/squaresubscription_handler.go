package handlersquare

// import (
// 	"context"
// 	"erp/api/dto"
// 	"erp/internal/app/config"
// 	dtosquare "erp/internal/app/plugin/square/api_square/dto_square"
// 	squareservice "erp/internal/app/plugin/square/square_service"
// 	"erp/internal/app/service/helpers"
// 	"fmt"
// 	"net/http"

// 	"github.com/danielgtaylor/huma/v2"
// )

// type SquareSubscriptionHandler struct {
// 	squareSubscriptionService *squareservice.SquareSubscriptionService
// 	sessionHelper             *helpers.SessionHelper
// 	locale                    helpers.Locale
// }

// func NewSquareSubscriptionHandler(
// 	api *huma.API,
// 	helpers *helpers.Helpers,
// 	squareSubscriptionService *squareservice.SquareSubscriptionService,
// 	base string,
// 	tag string,
// 	authMiddleware config.AppMiddleware,
// 	validateCompanyM config.AppMiddleware,
// ) {
	
// 	paths := NewSquareSubscriptionPaths(base)
// 	handler := SquareSubscriptionHandler{
// 		squareSubscriptionService: squareSubscriptionService,
// 		sessionHelper:             helpers.Session,
// 		locale:                    helpers.Locale,
// 	}
// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "cancel-square-subscription",
// 		Method:        http.MethodPost,
// 		Path:          paths.Cancel,
// 		Summary:       "Cancel square subscription",
// 		Tags:          []string{tag},
// 		DefaultStatus: http.StatusOK,
// 		Middlewares: huma.Middlewares{authMiddleware,validateCompanyM},
// 	}, handler.CancelSubscription)
// }

// func (h *SquareSubscriptionHandler) CancelSubscription(ctx context.Context, i *dtosquare.RequestSubscriptionCancel) (
// 	*dto.ResponseMessage, error,
// ) {
// 	var response dto.ResponseMessage
// 	req, _ := h.sessionHelper.GetSession(ctx)

// 	res, err := h.squareSubscriptionService.CancelSubscription(req, i)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest(h.locale.MustLocalize(
// 			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
// 			helpers.OptionsLocale.WithID("Error.FailToCancelSubscription"),
// 		), err)
// 	}
// 	fmt.Println("ERRORS",res)
// 	if res != nil {
// 		if len(res.Errors) > 0 {
// 			return nil, huma.Error400BadRequest(res.Errors[0].Detail, nil)
// 		}
// 	}
// 	response.Body.Message = h.locale.MustLocalize(
// 		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
// 		helpers.OptionsLocale.WithID("Message.CancelSubscription"),
// 	)
// 	return &response, nil
// }
