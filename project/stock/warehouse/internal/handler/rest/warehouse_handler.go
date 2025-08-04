package warehouse_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	warehouse_ucase "erp/project/stock/warehouse/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type WareHouseHandler struct {
	sessionHelper    helpers.SessionHelper
	locale           helpers.Locale
	warehouseUseCase warehouse_ucase.WarehouseUseCase
	errorHelper      helpers.ErrorHelper
	permission      repository.PermissionService
}

func NewWareHouseHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	warehouseUseCase warehouse_ucase.WarehouseUseCase,
	permission repository.PermissionService,
) {
	paths := NewWareHousePaths(domain.WAREHOUSE_BASE_ROUTE)
	tags :=  []string{"Warehouse"}
	h := WareHouseHandler{
		sessionHelper:    helpers.Session,
		locale:           helpers.Locale,
		warehouseUseCase: warehouseUseCase,
		errorHelper:      helpers.Error,
		permission:      permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "create-warehouse",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Warehouse",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateWareHouse)

	huma.Register(api, huma.Operation{
		OperationID:   "get-warehouses",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Warehouses",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetWareHouses)

	huma.Register(api, huma.Operation{
		OperationID:   "get-warehouse-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Warehouse Detail",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetWareHouseDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "get-warehouses-tree-view",
		Method:        http.MethodGet,
		Path:          paths.TreeView,
		Summary:       "Get Warehouses Tree View",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetTreeViewData)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-warehouse",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Warehouse",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditWarehouse)
}
func (h *WareHouseHandler) EditWarehouse(ctx context.Context, d *dto.EditWarehouseRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.warehouseUseCase.EditWarehouse(req, d)
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


func (h *WareHouseHandler) GetTreeViewData(ctx context.Context,i *struct{})(
	*dto.ResponseData[[]dto.TreeEntryDto],error){
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res,err := h.warehouseUseCase.GetWarehouseTreeView(req)
	if err != nil {
		return nil,h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
	}
	actions := h.permission.GetActions(req.Ctx,domain.WAREHOUSE.ID)
	var response dto.ResponseData[[]dto.TreeEntryDto]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response,nil
}

func (h *WareHouseHandler) GetWareHouseDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.WareHouseDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.warehouseUseCase.GetWareHouseDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.WAREHOUSE.ID)

	var response dto.EntityResponse[dto.ResultEntity[dto.WareHouseDto]]
	response.Body.Actions = actions
	response.Body.Result = res

	return &response, err
}

func (h *WareHouseHandler) GetWareHouses(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.WareHouseDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.warehouseUseCase.GetWareHouses(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
	}
	actions := h.permission.GetActions(req.Ctx,domain.WAREHOUSE.ID)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.WareHouseDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *WareHouseHandler) CreateWareHouse(ctx context.Context, i *dto.CreateWareHouseRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.warehouseUseCase.CreateWareHouse(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToCreateWareHouse"),
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
