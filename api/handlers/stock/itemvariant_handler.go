package stock

import (
	"context"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/internal/app/service/services/account_service"
	stockservice "erp/internal/app/service/services/stock_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ItemVaraintHandler struct {
	itemVaraintService *stockservice.ItemVariantService
	locale             helpers.Locale
	sessionHelper      helpers.SessionHelper
	errorHelper        helpers.ErrorHelper
	roleService        *account_service.RoleService
}

func NewItemVariantHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tags []string,
	middlewares huma.Middlewares,
) {
	paths := NewItemVariantPaths(base)
	handler := ItemVaraintHandler{
		itemVaraintService: services.ItemVariantService,
		locale:             helpers.Locale,
		sessionHelper:      helpers.Session,
		roleService:        services.RoleService,
		errorHelper: helpers.Error,
	}

	huma.Register(*api, huma.Operation{
		OperationID:   "create-item-variant",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create item variant",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateItemVariant)

	huma.Register(*api, huma.Operation{
		OperationID:   "get-variant-from-item",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get variant from items",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetVariantsFromItem)

}

func (h *ItemVaraintHandler) GetVariantsFromItem(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.ItemVariantDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.itemVaraintService.GetVariantsFromItem(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
	}
	actions := h.roleService.GetActions(req,domain.ITEM)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.ItemVariantDto]]{}
	response.Body.Actions = actions
	response.Body.PaginationResult = res
	return response, err

}

func (h *ItemVaraintHandler) CreateItemVariant(ctx context.Context, i *dto.CreateItemVariantRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.itemVaraintService.CreateItemVariant(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err,"Error.FailToCreateItemVariant")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedItemVariantSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}
