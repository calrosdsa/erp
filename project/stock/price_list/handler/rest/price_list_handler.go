package rest_price_list

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	price_list_ucase "erp/project/stock/price_list/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PriceListHandler struct {
	priceListUcase price_list_ucase.PriceListUseCase
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	permission    repository.PermissionService
}

func NewPriceListHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	priceListUcase price_list_ucase.PriceListUseCase,
	permission repository.PermissionService,
) {
	paths := NewPriceListPaths(domain.PRICE_LIST_BASE_ROUTE)
	tags := []string{"Price List"}
	handler := PriceListHandler{
		priceListUcase: priceListUcase,
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
		permission:    permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "get-price-list-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Retrieve Price List Detail",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPriceListDetail)
	huma.Register(api, huma.Operation{
		OperationID:   "get-price-lists",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Retrieve Price List items",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPriceLists)

	huma.Register(api, huma.Operation{
		OperationID:   "create-price-list",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create price list",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CreatePriceList)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-price-list",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Price List",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditPriceList)
}
func (h *PriceListHandler) EditPriceList(ctx context.Context, d *dto.EditPriceListRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.priceListUcase.EditPriceList(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}


func (h *PriceListHandler) CreatePriceList(ctx context.Context, i *dto.CreatePriceListRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	_,err := h.priceListUcase.CreatePriceList(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreatePriceList")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedPriceListSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *PriceListHandler) GetPriceListDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.PriceListDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.priceListUcase.GetListPriceDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.PRICE_LIST.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.PriceListDto]]
	response.Body.Actions = actions
	response.Body.Result = res
	return &response, nil
}
func (h *PriceListHandler) GetPriceLists(ctx context.Context, i *dto.RequestPriceLists) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.PriceListDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.priceListUcase.GetPriceLists(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.PRICE_LIST.ID)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.PriceListDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

// func (h *PriceListHandler) UpdateItemPriceList(ctx context.Context, i *dto.UpsertPriceListRequest) (*dto.ResponseMessage, error) {
// 	req, _ := h.sessionHelper.GetSession(ctx)
// 	var response dto.ResponseMessage
// 	err := h.priceListUcase.UpsertPriceList(req, i)
// 	// h.locale.MustLocalize()
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Fail to upsert price list")
// 	}
// 	return &response, err
// }

// func (h *PriceListHandler) CreateItemPriceList(ctx context.Context, i *dto.UpsertPriceListRequest) (*dto.ResponseMessage, error) {
// 	req, _ := h.sessionHelper.GetSession(ctx)
// 	var response dto.ResponseMessage
// 	err := h.priceListUcase.UpsertPriceList(req, i)
// 	// h.locale.MustLocalize()
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Fail to upsert price list")
// 	}
// 	return &response, err
// }
