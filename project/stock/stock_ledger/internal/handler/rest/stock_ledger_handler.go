package stock_ledger_rest

import (
	"context"

	// "erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	stock_ledger_ucase "erp/project/stock/stock_ledger/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type StockLedgerHandler struct {
	stockLedgerUcase  stock_ledger_ucase.StockLedgerUseCase
	sessionHelper    helpers.SessionHelper
	errorHelper      helpers.ErrorHelper
	locale           helpers.Locale
	permissioService repository.PermissionService
}

func NewStockLedgerHandler(
	api huma.API,
	middlewares huma.Middlewares,
	helpers *helpers.Helpers,
	stockLedgerUcase stock_ledger_ucase.StockLedgerUseCase,
	permissionService repository.PermissionService,
) {
	paths := NewStockLedgerPaths(domain.STOCK_LEDGER_BASE_ROUTE)
	tag := []string{"Stock Ledger Report"}
	handler := StockLedgerHandler{
		stockLedgerUcase:  stockLedgerUcase,
		sessionHelper:    helpers.Session,
		errorHelper:      helpers.Error,
		locale:           helpers.Locale,
		permissioService: permissionService,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "get-stock_ledger-report",
		Method:        http.MethodGet,
		Path:          paths.StockLedgerReport,
		Summary:       "Get Stock Ledger Report",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetStockLedgerReport)

	huma.Register(api, huma.Operation{
		OperationID:   "get-stock-balance-report",
		Method:        http.MethodGet,
		Path:          paths.StockBalanceReport,
		Summary:       "Get Stock Balance Report",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetStockBalanceReport)
}

func (h *StockLedgerHandler) GetStockBalanceReport(ctx context.Context, i *dto.RequestStockBalance) (
	*dto.EntityResponse[[]dto.StockBalanceEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.stockLedgerUcase.GetStockBalanceReport(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.StockBalanceEntryDto]
	response.Body.Result = res
	return &response, nil
}


func (h *StockLedgerHandler) GetStockLedgerReport(ctx context.Context, i *dto.RequestStockLedger) (
	*dto.EntityResponse[[]dto.StockLedgerEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.stockLedgerUcase.GetStockLedgerReport(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.StockLedgerEntryDto]
	response.Body.Result = res
	return &response, nil
}
