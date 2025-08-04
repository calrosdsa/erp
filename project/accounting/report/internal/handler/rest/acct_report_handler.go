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

type AcctReportHandler struct {
	acctReportUcase  acct_report_ucase.AcctReportUseCase
	sessionHelper    helpers.SessionHelper
	errorHelper      helpers.ErrorHelper
	locale           helpers.Locale
	permissioService repository.PermissionService
}

func NewAcctReportHandler(
	api huma.API,
	middlewares huma.Middlewares,
	helpers *helpers.Helpers,
	acctReportUcase acct_report_ucase.AcctReportUseCase,
	permissionService repository.PermissionService,
) {
	paths := NewAcctReportPaths(domain.ACCT_REPORT_BASE_ROUTE)
	tag := []string{"Acct Report"}
	handler := AcctReportHandler{
		acctReportUcase:  acctReportUcase,
		sessionHelper:    helpers.Session,
		errorHelper:      helpers.Error,
		locale:           helpers.Locale,
		permissioService: permissionService,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "get-general-ledger",
		Method:        http.MethodGet,
		Path:          paths.GeneralLedger,
		Summary:       "Get General Ledger",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetGeneralLedger)

	huma.Register(api, huma.Operation{
		OperationID:   "get-account-payable",
		Method:        http.MethodGet,
		Path:          paths.AccountPayable,
		Summary:       "Get Account Payable",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAccountPayable)

	huma.Register(api, huma.Operation{
		OperationID:   "get-account-payable-sumary",
		Method:        http.MethodGet,
		Path:          paths.AccountPayableSumary,
		Summary:       "Get Account Payable Sumary",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAccountPayableSumary)

	huma.Register(api, huma.Operation{
		OperationID:   "account-receivable",
		Method:        http.MethodGet,
		Path:          paths.AccountReceivable,
		Summary:       "Account Receivable",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAccountReceivable)

	huma.Register(api, huma.Operation{
		OperationID:   "receivable-sumary",
		Method:        http.MethodGet,
		Path:          paths.AccountReceivableSumary,
		Summary:       "Receivable Sumary",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAccountReceivableSumary)

	huma.Register(api, huma.Operation{
		OperationID:   "account-balance",
		Method:        http.MethodGet,
		Path:          paths.AccountBalance,
		Summary:       "Account Balance",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAccountBalance)
}
func (h *AcctReportHandler) GetAccountBalance(ctx context.Context, i *dto.RequestAccountBalance) (
	*dto.EntityResponse[dto.GeneralLedgerOpening], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.acctReportUcase.GetAccountBalance(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.GeneralLedgerOpening]
	response.Body.Result = res
	return &response, nil
}

func (h *AcctReportHandler) GetAccountReceivable(ctx context.Context, i *dto.RequestAccountReceivable) (
	*dto.EntityResponse[[]dto.AccountReceivableEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.acctReportUcase.GetAccountReceivable(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.AccountReceivableEntryDto]
	response.Body.Result = res
	return &response, nil
}

func (h *AcctReportHandler) GetAccountReceivableSumary(ctx context.Context, i *dto.RequestAccountReceivable) (
	*dto.EntityResponse[[]dto.SumaryEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.acctReportUcase.GetAccountReceivableSumary(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissioService.GetActions(req.Ctx, domain.LEDGER.ID)
	var response dto.EntityResponse[[]dto.SumaryEntryDto]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *AcctReportHandler) GetAccountPayableSumary(ctx context.Context, i *dto.RequestAccountPayable) (
	*dto.EntityResponse[[]dto.SumaryEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.acctReportUcase.GetAccountPayableSumary(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.SumaryEntryDto]
	response.Body.Result = res
	return &response, nil
}

func (h *AcctReportHandler) GetAccountPayable(ctx context.Context, i *dto.RequestAccountPayable) (
	*dto.EntityResponse[[]dto.AccountPayableEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.acctReportUcase.GetAccountPayable(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissioService.GetActions(req.Ctx, domain.LEDGER.ID)
	var response dto.EntityResponse[[]dto.AccountPayableEntryDto]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *AcctReportHandler) GetGeneralLedger(ctx context.Context, i *dto.RequestGeneralLedger) (
	*dto.EntityResponse[dto.GeneralLedgerData], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.acctReportUcase.GetGeneralLedgerReport(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissioService.GetActions(req.Ctx, domain.LEDGER.ID)
	var response dto.EntityResponse[dto.GeneralLedgerData]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}
