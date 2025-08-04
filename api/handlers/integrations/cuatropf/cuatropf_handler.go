package cuatropf

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	clientservice "erp/internal/app/service/services/client_service"
	"erp/internal/app/service/services/company_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type CuatropfHandler struct {
	companyService *company_service.CompanyService
	sessionHelper  helpers.SessionHelper
	clientService  *clientservice.ClientService
	locale         helpers.Locale
}

func NewHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag string,
	// middlewares huma.Middlewares,
) {
	paths := NewCuatropfPath(base)
	handler := CuatropfHandler{
		companyService: services.CompanyService,
		sessionHelper:  helpers.Session,
		clientService:  services.ClientService,
		locale:         helpers.Locale,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "cuatropf-subscription",
		Method:        http.MethodPost,
		Path:          paths.Subscription,
		Summary:       "Cuatropf Subscription",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		// Middlewares:   middlewares,
	}, handler.Subscription)
}

func (h *CuatropfHandler) Subscription(ctx context.Context, i *dto.CuatropfSubscriptionRequest) (
	*dto.ResponseMessage, error,
) {
	var response dto.ResponseMessage
	var req common.RequestContext
	// company, err := h.companyService.GetCompanyByUuid(ctx, i.CompanyUuid)
	// if err != nil {
	// 	return nil, huma.Error400BadRequest("No company Found", err)
	// }
	req.Ctx = ctx
	// req.ActiveCompany = company
	var clientRequest dto.CreateClientRequest
	clientRequest.Body.ClientRequestDto = i.Body.ClientRequestDto
	// _, err = h.clientService.CreateCustomer(&req, &clientRequest)
	// if err != nil {
	// 	return nil, huma.Error400BadRequest(
	// 		h.locale.MustLocalize(
	// 			helpers.OptionsLocale.WithID("Error.FailToCancelSubscription"),
	// 			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	// 		),
	// 		err)
	// }
	return &response, nil
}

// func (h CuatropfHandler) SignIn(ctx echo.Context) error {
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
