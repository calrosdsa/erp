package item_inventory_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	item_inventory_ucase "erp/project/stock/item_inventory/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type itemInventoryHandler struct {
	sessionHelper      helpers.SessionHelper
	locale             helpers.Locale
	errorHelper        helpers.ErrorHelper
	itemInventoryUcase item_inventory_ucase.ItemInventoryUcase
	permission         repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	itemInventoryUcase item_inventory_ucase.ItemInventoryUcase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.ITEM_BASE_ROUTE
	tags := []string{"Item Inventory"}
	path := NewPaths(base)
	h := itemInventoryHandler{
		sessionHelper:      helpers.Session,
		locale:             helpers.Locale,
		errorHelper:        helpers.Error,
		itemInventoryUcase: itemInventoryUcase,
		permission:         permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "item-inventory-setting",
		Method:        http.MethodGet,
		Summary:       "item-inventory-setting",
		Tags:          tags,
		Path:          path.InventorySettingsDetail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetItemInventory)
	huma.Register(api, huma.Operation{
		OperationID:   "edit-item-inventory-setting",
		Method:        http.MethodPut,
		Summary:       "edit-item-inventory-setting",
		Tags:          tags,
		Path:          path.InventorySettings,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditItemInventory)
}

func (h *itemInventoryHandler) EditItemInventory(ctx context.Context, d *dto.EditItemInventoryRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.itemInventoryUcase.EditItemInventory(req, d.Body)
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

func (h *itemInventoryHandler) GetItemInventory(ctx context.Context, d *dto.RequestEntity) (
	*dto.ResponseData[dto.ItemInventoryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.itemInventoryUcase.GetItemInventory(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[dto.ItemInventoryDto]
	response.Body.Result = res
	return &response, nil
}
