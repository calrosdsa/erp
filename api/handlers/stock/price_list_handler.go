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

type ItemPriceListHandler struct {
	priceListService *stockservice.PriceListService
	sessionHelper    helpers.SessionHelper
	locale           helpers.Locale
	errorHelper      helpers.ErrorHelper
	roleService *account_service.RoleService
}

func NewPriceListHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tags []string,
	middlewares huma.Middlewares,
) {
	paths := NewPriceListPaths(base)
	handler := ItemPriceListHandler{
		priceListService: services.PriceListService,
		sessionHelper:    helpers.Session,
		locale:           helpers.Locale,
		errorHelper:      helpers.Error,
		roleService: services.RoleService,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "get-price-list-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Retrieve Price List Detail",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPriceListDetail)
	huma.Register(*api, huma.Operation{
		OperationID:   "get-price-lists",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Retrieve Price List items",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPriceLists)

	huma.Register(*api, huma.Operation{
		OperationID:   "create-price-list",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create price list",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CreatePriceList)

	huma.Register(*api, huma.Operation{
		OperationID:   "update-item-price-list",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Update item price list",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateItemPriceList)
}

func (h *ItemPriceListHandler) CreatePriceList(ctx context.Context, i *dto.CreatePriceListRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.priceListService.CreatePriceList(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err,"Error.FailToCreatePriceList")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedPriceListSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ItemPriceListHandler) GetPriceListDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.PriceListDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.priceListService.GetListPriceDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
	}
	actions := h.roleService.GetActions(req,domain.PRICE_LIST)
	var response dto.EntityResponse[dto.ResultEntity[dto.PriceListDto]]
	response.Body.Actions =actions
	response.Body.Result = res
	return &response, nil
}
func (h *ItemPriceListHandler) GetPriceLists(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.PriceListDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.priceListService.GetPriceLists(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
	}
	actions := h.roleService.GetActions(req,domain.PRICE_LIST)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.PriceListDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *ItemPriceListHandler) UpdateItemPriceList(ctx context.Context, i *dto.UpsertPriceListRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	var response dto.ResponseMessage
	err := h.priceListService.UpsertPriceList(req, i)
	// h.locale.MustLocalize()
	if err != nil {
		return nil, huma.Error400BadRequest("Fail to upsert price list")
	}
	return &response, err
}

// func (h *ItemPriceListHandler) CreateItemPriceList(ctx context.Context, i *dto.UpsertPriceListRequest) (*dto.ResponseMessage, error) {
// 	req, _ := h.sessionHelper.GetSession(ctx)
// 	var response dto.ResponseMessage
// 	err := h.priceListService.UpsertPriceList(req, i)
// 	// h.locale.MustLocalize()
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Fail to upsert price list")
// 	}
// 	return &response, err
// }
