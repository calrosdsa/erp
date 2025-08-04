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

// optional code omitted

type ItemAttributeHandler struct {
	sessionHelper        helpers.SessionHelper
	itemAttributeService *stockservice.ItemAttributeService
	locale               helpers.Locale
	roleService          *account_service.RoleService
	errorHelper          helpers.ErrorHelper
}

func NewItemAttributeHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tags []string,
	middlewares huma.Middlewares,
) {
	paths := NewItemAttributePaths(base)
	handler := ItemAttributeHandler{
		sessionHelper:        helpers.Session,
		itemAttributeService: services.ItemAttributeService,
		locale:               helpers.Locale,
		roleService:          services.RoleService,
		errorHelper:          helpers.Error,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "create-item-attribute",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create item attribute",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CreateItemAttribute)
	huma.Register(*api, huma.Operation{
		OperationID:   "item-attributes",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Item attributes",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemAttributes)

	huma.Register(*api, huma.Operation{
		OperationID:   "item-attribute",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Item attribute",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemAttribute)

	huma.Register(*api, huma.Operation{
		OperationID:   "update-item-attribute-value",
		Method:        http.MethodPut,
		Path:          paths.ItemAttributeValue,
		Summary:       "Update Item Attribute Value",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpsertItemAttributeValue)

	huma.Register(*api, huma.Operation{
		OperationID:   "create-item-attribute-value",
		Method:        http.MethodPost,
		Path:          paths.ItemAttributeValue,
		Summary:       "Create Item Attribute Value",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.UpsertItemAttributeValue)

}

func (h *ItemAttributeHandler) UpsertItemAttributeValue(ctx context.Context, i *dto.UpsertItemAttributeValueRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.itemAttributeService.UpsertItemAttributeValue(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateItemAttributeValue")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateItemAttributeValueSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ItemAttributeHandler) GetItemAttribute(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ItemAttributeDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	req.Params = make(map[string]string)
	h.sessionHelper.AppendParam(req, helpers.COLUMN_PARAM, i.OrderColumn)
	h.sessionHelper.AppendParam(req, helpers.ORDER_PARAM, i.Order)
	res, err := h.itemAttributeService.GetItemAttributeDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.roleService.GetActions(req, domain.ITEM_ATTRIBUTE)
	var response dto.EntityResponse[dto.ResultEntity[dto.ItemAttributeDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil

}

func (h *ItemAttributeHandler) CreateItemAttribute(ctx context.Context, i *dto.CreateItemAttributeRequest) (
	*dto.EntityResponse[dto.ItemAttributeDto], error) {
	req, _ := h.sessionHelper.GetSession(ctx)

	res, err := h.itemAttributeService.CreateAttributeItem(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateItemAttribute")
	}
	response := dto.EntityResponse[dto.ItemAttributeDto]{}
	response.Body.Result = res
	return &response, nil
}

func (h *ItemAttributeHandler) GetItemAttributes(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.ItemAttributeDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.itemAttributeService.GetItemAttributes(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.roleService.GetActions(req, domain.ITEM_ATTRIBUTE)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.ItemAttributeDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}
