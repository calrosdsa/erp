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

type ItemStockHandler struct {
	sessionHelper    helpers.SessionHelper
	itemStockService *stockservice.ItemStockService
	locale           helpers.Locale
	errorHelper      helpers.ErrorHelper
	roleService      *account_service.RoleService
}

func NewItemStockHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag []string,
	middlewares huma.Middlewares,
) {
	paths := NewItemStockPaths(base)
	handler := ItemStockHandler{
		sessionHelper:    helpers.Session,
		itemStockService: services.ItemStockService,
		locale:           helpers.Locale,
		roleService:      services.RoleService,
		errorHelper:      helpers.Error,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "add-item-to-warehouse",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Add item to warehouse",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.AddItemToWareHouse)

	huma.Register(*api, huma.Operation{
		OperationID:   "edit-item-to-warehouse",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit item to warehouse",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditItemToWareHouse)

	huma.Register(*api, huma.Operation{
		OperationID:   "get-item-stock-levels",
		Method:        http.MethodGet,
		Path:          paths.Item,
		Summary:       "Retrieve item stock levels",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemStockLevels)

	huma.Register(*api, huma.Operation{
		OperationID:   "get-warehouse-stock-levels",
		Method:        http.MethodGet,
		Path:          paths.Warehouse,
		Summary:       "Retrieve warehouse stock levels",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetWareHouseItemStockLevels)
}

func (h *ItemStockHandler) GetWareHouseItemStockLevels(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.StockLevelDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.itemStockService.GetWarehouseItemStockLevels(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest("Error", err)
	}
	actions := h.roleService.GetActions(req,domain.ITEM_STOCK)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.StockLevelDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions

	return response, err
}

func (h *ItemStockHandler) GetItemStockLevels(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.StockLevelDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.itemStockService.GetItemStockLevels(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.roleService.GetActions(req, domain.ITEM_STOCK)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.StockLevelDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *ItemStockHandler) AddItemToWareHouse(ctx context.Context, i *dto.AddStockLevelRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.itemStockService.AddStockLevel(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToAddItemStock"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateWareHouseSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ItemStockHandler) EditItemToWareHouse(ctx context.Context, i *dto.AddStockLevelRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.itemStockService.AddStockLevel(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToEditItemStock"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditItemStockSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

// func (h *ItemStockHandler) GetUnitOfMeasures(ctx context.Context, i *dto.UOMsRequest) (*dto.UOMsResponse, error) {
// 	var response dto.UOMsResponse
// 	req, err := h.sessionHelper.GetSession(ctx)
// 	if err != nil {
// 		return nil, huma.Error401Unauthorized("Not Authorized", err)
// 	}
// 	res,err := h.uomService.GetUnitOfMeasures(req,i)
// 	if err != nil {
// 		return nil,huma.NewError(http.StatusBadRequest,"Failed to retrieve UOM (Unit of Measure)")
// 	}
// 	response.Body.Results = res
// 	return &response, nil
// }
