package rest_stock_entry

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	stock_entry_ucase "erp/project/stock/stock_entry/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type StockEntryHandler struct {
	sessionHelper   helpers.SessionHelper
	locale          helpers.Locale
	errorHelper     helpers.ErrorHelper
	stockEntryUcase stock_entry_ucase.StockEntryUseCase
	permission      repository.PermissionService
}

func NewStockEntryHandler(
	api huma.API,
	helpers *helpers.Helpers,
	stockEntryUcase stock_entry_ucase.StockEntryUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.STOCK_ENTRY_BASE_ROUTE
	tags := []string{"Stock Entry"}
	path := NewStockEntryPaths(base)
	h := StockEntryHandler{
		sessionHelper:   helpers.Session,
		locale:          helpers.Locale,
		errorHelper:     helpers.Error,
		stockEntryUcase: stockEntryUcase,
		permission:      permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "stock-entries",
		Method:        http.MethodGet,
		Summary:       "Stock Entries",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetStockEntries)
	huma.Register(api, huma.Operation{
		OperationID:   "stock-entry",
		Method:        http.MethodGet,
		Summary:       "Stock Entry",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetStockEntry)

	huma.Register(api, huma.Operation{
		OperationID:   "create-stock-entry",
		Method:        http.MethodPost,
		Summary:       "Create Stock Entry",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateStockEntry)

	huma.Register(api, huma.Operation{
		OperationID:   "update-stock-entry-status",
		Method:        http.MethodPut,
		Summary:       "Update Stock Entry Status",
		Path:          path.UpdateStatus,
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStockEntryStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-stock-entry",
		Method:        http.MethodPut,
		Summary:       "Edit Stock Entry",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditStockEntry)
}
func (h *StockEntryHandler) EditStockEntry(ctx context.Context, d *dto.EditStockEntryRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.stockEntryUcase.EditStockEntry(req, d.Body)
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

func (h *StockEntryHandler) UpdateStockEntryStatus(ctx context.Context, i *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.stockEntryUcase.UpdateStockEntryStatus(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdate")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}
func (h *StockEntryHandler) CreateStockEntry(ctx context.Context, d *dto.CreateStockEntryRequest) (
	*dto.ResponseData[dto.StockEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.stockEntryUcase.CreateStockEntry(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseData[dto.StockEntryDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *StockEntryHandler) GetStockEntry(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.StockEntryDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.stockEntryUcase.GetStockEntry(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.StockEntryDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.STOCK_ENTRY.ID)
	response.Body.AssociatedActions = h.getExtraActions(req.Ctx)
	return &response, nil
}

func (h *StockEntryHandler) GetStockEntries(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.StockEntryDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.stockEntryUcase.GetStockEntries(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.StockEntryDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.STOCK_ENTRY.ID)
	return &response, nil
}

func (h *StockEntryHandler) getExtraActions(ctx context.Context) map[int][]dto.ActionDto {
	r := make(map[int][]dto.ActionDto)
	r[int(domain.GENERAL_LEDGER.ID)] = h.permission.GetActions(ctx, domain.GENERAL_LEDGER.ID)
	return r
}
