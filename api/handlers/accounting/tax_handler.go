package accounting

import (
	"context"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/internal/app/service/services/account_service"
	accountingservice "erp/internal/app/service/services/accounting_service"

	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type TaxHandler struct {
	sessionHelper helpers.SessionHelper
	taxService    *accountingservice.TaxService
	locale        helpers.Locale
	roleService   *account_service.RoleService
	errorHelper   helpers.ErrorHelper
}

func NewTaxHandler(api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tags []string,
	commonMiddlewares huma.Middlewares,
) {

	paths := NewTaxPaths(base)
	handler := TaxHandler{
		sessionHelper: helpers.Session,
		taxService:    services.TaxService,
		locale:        helpers.Locale,
		roleService:   services.RoleService,
		errorHelper:   helpers.Error,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "get-tax-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Tax Detail",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   commonMiddlewares,
	}, handler.GetTaxDetail)

	huma.Register(*api, huma.Operation{
		OperationID:   "get-taxes",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Taxes",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   commonMiddlewares,
	}, handler.GetTaxes)

	huma.Register(*api, huma.Operation{
		OperationID:   "create-tax",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Tax",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   commonMiddlewares,
	}, handler.CreateTax)
	// huma.Register(*api,huma.Operation{
	// 	OperationID: "edit-tax",
	// })
}

func (h *TaxHandler) CreateTax(ctx context.Context, i *dto.CreateTaxRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.taxService.CreateTax(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToCreateTax"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateTaxSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *TaxHandler) GetTaxes(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.TaxDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.taxService.GetTaxes(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.roleService.GetActions(req, domain.TAX)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.TaxDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *TaxHandler) GetTaxDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.TaxDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.taxService.GetTaxDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.roleService.GetActions(req, domain.TAX)
	var response dto.EntityResponse[dto.ResultEntity[dto.TaxDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, err
}

// func (h TaxHandler) SignIn(ctx echo.Context) error {
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
