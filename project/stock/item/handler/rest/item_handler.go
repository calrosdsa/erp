package item_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	item_ucase "erp/project/stock/item/usecase"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ItemHandler struct {
	itemUseCase   item_ucase.ItemUseCase
	permission    repository.PermissionService
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
}

func NewItemHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	itemUseCase item_ucase.ItemUseCase,
	permission repository.PermissionService,
) {
	paths := NewItemPath(domain.ITEM_BASE_ROUTE)
	tags := []string{"Item"}
	handler := ItemHandler{
		itemUseCase:   itemUseCase,
		permission:    permission,
		sessionHelper: helpers.Session,
		errorHelper:   helpers.Error,
		locale:        helpers.Locale,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "get-items",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Items",
		Description:   "Retrieve Items",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.GetItems)
	huma.Register(api, huma.Operation{
		OperationID:   "get-item-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Retrieve item detail base on code",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.GetItemDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "create-item",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create item",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateItem)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-item",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit item",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.EditItem)

	huma.Register(api, huma.Operation{
		OperationID:   "item-actions",
		Method:        http.MethodGet,
		Path:          paths.Actions,
		Summary:       "Get Item Actions",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemActions)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-item",
		Method:        http.MethodPut,
		Path:          paths.Base + "/update-status",
		Summary:       "Update Status Item",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateStatus)
}

func (m *ItemHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (_a0 *dto.ResponseMessage, _a1 error) {

	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.itemUseCase.UpdateStatus(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ItemHandler) GetItemActions(ctx context.Context, i *struct{}) (
	*dto.ResponseData[any], error,
) {
	_, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	var res dto.ResponseData[any]
	res.Body.Actions = h.permission.GetActions(ctx, domain.ITEM.ID)
	res.Body.AssociatedActions = h.getExtraActions(ctx)
	return &res, nil
}

func (h *ItemHandler) GetItemDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ItemDetailDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.itemUseCase.GetItemDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.ITEM.ID)

	response := &dto.EntityResponse[dto.ResultEntity[dto.ItemDetailDto]]{}
	response.Body.Result = res
	response.Body.Actions = actions
	response.Body.AssociatedActions = h.getExtraActions(ctx)
	return response, err
}

func (h *ItemHandler) GetItems(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.ItemDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.itemUseCase.GetItems(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest("Error", err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.ITEM.ID)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.ItemDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *ItemHandler) EditItem(ctx context.Context, d *dto.ItemRequestData) (*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.itemUseCase.EditItem(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ItemHandler) CreateItem(ctx context.Context, i *dto.ItemRequestData) (*dto.EntityResponse[dto.ItemDto], error) {
	fmt.Println("CREATE ITEM")
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	fmt.Println("CREATE ITEM2")
	res, err := h.itemUseCase.CreateItem(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateItem")
	}
	var response dto.EntityResponse[dto.ItemDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, err
}

func (h *ItemHandler) getExtraActions(ctx context.Context) map[int][]dto.ActionDto {
	r := make(map[int][]dto.ActionDto)
	r[int(domain.ITEM_PRICE.ID)] = h.permission.GetActions(ctx, domain.ITEM_PRICE.ID)
	r[int(domain.ITEM_GROUP.ID)] = h.permission.GetActions(ctx, domain.ITEM_GROUP.ID)
	return r
}
