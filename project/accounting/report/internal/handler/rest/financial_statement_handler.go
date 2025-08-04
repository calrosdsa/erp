package acct_report_rest

import (
	"context"

	// "erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	acct_report_ucase "erp/project/accounting/report/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type FinancialStatementHandler struct {
	financialStatementUcase acct_report_ucase.FinancialStatementUcase
	sessionHelper           helpers.SessionHelper
	errorHelper             helpers.ErrorHelper
	locale                  helpers.Locale
	permissioService        repository.PermissionService
}

func NewFinancialStatementHandler(
	api huma.API,
	middlewares huma.Middlewares,
	helpers *helpers.Helpers,
	financialStatementUcase acct_report_ucase.FinancialStatementUcase,
	permissionService repository.PermissionService,
) {
	paths := NewFinancialStatementPaths(domain.FINANCIAL_STATEMENT_BASE_ROUTE)
	tag := []string{"Financial Statement"}
	handler := FinancialStatementHandler{
		financialStatementUcase: financialStatementUcase,
		sessionHelper:           helpers.Session,
		errorHelper:             helpers.Error,
		locale:                  helpers.Locale,
		permissioService:        permissionService,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "profit-and-loss",
		Method:        http.MethodGet,
		Path:          paths.ProfitAndLoss,
		Summary:       "Profit and Loss Statement",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetProfitAndLostStatement)

	huma.Register(api, huma.Operation{
		OperationID:   "cash-flow",
		Method:        http.MethodGet,
		Path:          paths.CashFlow,
		Summary:       "Cash flow Statement",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CashFlowStatement)

	huma.Register(api, huma.Operation{
		OperationID:   "balance-sheet",
		Method:        http.MethodGet,
		Path:          paths.BalanceSheet,
		Summary:       "Balance Sheet",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.BalanceSheetStatement)
}

func (h *FinancialStatementHandler) BalanceSheetStatement(ctx context.Context, i *dto.RequestFinancialStatement) (
	*dto.EntityResponse[[]dto.BalanceSheetEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.financialStatementUcase.BalanceSheetStatement(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.BalanceSheetEntryDto]
	response.Body.Result = res
	return &response, nil
}

func (h *FinancialStatementHandler) CashFlowStatement(ctx context.Context, i *dto.RequestFinancialStatement) (
	*dto.EntityResponse[[]dto.CashFlowEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.financialStatementUcase.CashFlowStatement(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.CashFlowEntryDto]
	response.Body.Result = res
	return &response, nil
}


func (h *FinancialStatementHandler) GetProfitAndLostStatement(ctx context.Context, i *dto.RequestFinancialStatement) (
	*dto.EntityResponse[[]dto.ProfitAndLossEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.financialStatementUcase.ProfitAndLossStatement(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.ProfitAndLossEntryDto]
	response.Body.Result = res
	return &response, nil
}
