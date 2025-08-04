package itemline_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	itemline_ucase "erp/project/stock/itemline/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ItemLineHandler struct {
	sessionHelper   helpers.SessionHelper
	permission      repository.PermissionService
	errorHelper     helpers.ErrorHelper
	itemlineUseCase itemline_ucase.ItemLineUseCase
	locale          helpers.Locale
}

func NewItemLineHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	itemLineUseCase itemline_ucase.ItemLineUseCase,
) {
	paths := NewItemLinePaths(domain.ITEM_LINE_BASE_ROUTE)
	tag := []string{"Item Line"}
	handler := ItemLineHandler{
		sessionHelper:   helpers.Session,
		permission:      permission,
		itemlineUseCase: itemLineUseCase,
		errorHelper:     helpers.Error,
		locale:          helpers.Locale,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "update-item-line",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "EditItemLine",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditItemLine)

	huma.Register(api, huma.Operation{
		OperationID:   "item-lines",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Item Lines",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemLines)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-line-item",
		Method:        http.MethodDelete,
		Path:          paths.Base,
		Summary:       "Delete Line Item",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.DeleteLineItem)

	huma.Register(api, huma.Operation{
		OperationID:   "add-line-item",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Add Line Item",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.AddLineItem)

	huma.Register(api, huma.Operation{
		OperationID:   "upsert-product-list",
		Method:        "POST",
		Path:          paths.ProductList,
		Summary:       "Upsert Product List",
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpsertProductList)
}

func (h *ItemLineHandler) UpsertProductList(ctx context.Context, d *dto.ProductListDataRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.itemlineUseCase.UpsertProductList(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")	
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.SuccessfullyMessage"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ItemLineHandler) AddLineItem(ctx context.Context, i *dto.AddLineItemRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.itemlineUseCase.AddLineItem(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ItemLineHandler) DeleteLineItem(ctx context.Context, i *dto.DeleteLineItemRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.itemlineUseCase.DeleteLineItem(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.DeletedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ItemLineHandler) GetItemLines(ctx context.Context, i *dto.RequestLineItems) (
	*dto.ResponseData[[]dto.LineItemDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.itemlineUseCase.GetItemLines(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.LineItemDto]
	response.Body.Result = res
	return &response, nil
}

func (h *ItemLineHandler) EditItemLine(ctx context.Context, i *dto.EditLineItemRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.itemlineUseCase.EditItemLine(req, i)
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
