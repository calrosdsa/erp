package ledger_rest

import (
	"context"

	// "erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	ledger_ucase "erp/project/accounting/ledger/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type LedgerHandler struct {
	ledgerUseCase    ledger_ucase.LedgerUseCase
	sessionHelper    helpers.SessionHelper
	errorHelper      helpers.ErrorHelper
	locale           helpers.Locale
	permissioService repository.PermissionService
}

func NewLedgerHandler(
	api huma.API,
	middlewares huma.Middlewares,
	helpers *helpers.Helpers,
	ledgerUseCase ledger_ucase.LedgerUseCase,
	permissionService repository.PermissionService,
) {
	paths := NewLedgerPaths(domain.LEDGER_BASE_ROUTE)
	tags := []string{"Ledger"}
	h := LedgerHandler{
		ledgerUseCase:    ledgerUseCase,
		sessionHelper:    helpers.Session,
		errorHelper:      helpers.Error,
		locale:           helpers.Locale,
		permissioService: permissionService,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "create-ledger",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Ledger",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.CreateLedger)

	huma.Register(api, huma.Operation{
		OperationID:   "get-acconts",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Accounts",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetLedgersAccount)

	huma.Register(api, huma.Operation{
		OperationID:   "get-ledger-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Ledger Detail",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetLedgerDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "get-general-ledger",
		Method:        http.MethodGet,
		Path:          paths.GeneralLedger,
		Summary:       "Get General Ledger",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetGeneralLedger)

	huma.Register(api, huma.Operation{
		OperationID:   "get-ledgers-tree-view",
		Method:        http.MethodGet,
		Path:          paths.TreeView,
		Summary:       "Get Ledgers Tree View",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetLedgersTreeView)
	huma.Register(api, huma.Operation{
		OperationID:   "edit-ledger",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Ledger",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditLedger)
}
func (h *LedgerHandler) EditLedger(ctx context.Context, d *dto.LedgerDataRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.ledgerUseCase.EditLedger(req, d.Body)
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

func (h *LedgerHandler) GetLedgersTreeView(ctx context.Context, i *struct{}) (
	*dto.ResponseData[[]dto.TreeEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.ledgerUseCase.GetLedgerAccountsTree(req)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissioService.GetActions(req.Ctx, domain.LEDGER.ID)
	var response dto.ResponseData[[]dto.TreeEntryDto]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *LedgerHandler) GetGeneralLedger(ctx context.Context, i *dto.RequestGeneralLedger) (
	*dto.EntityResponse[[]dto.GeneralLedgerEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.ledgerUseCase.GetGeneralLedgerReport(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissioService.GetActions(req.Ctx, domain.LEDGER.ID)
	var response dto.EntityResponse[[]dto.GeneralLedgerEntryDto]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *LedgerHandler) GetLedgerDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.LedgerDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.ledgerUseCase.GetLedgerDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissioService.GetActions(req.Ctx, domain.LEDGER.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.LedgerDetailDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}
func (h *LedgerHandler) GetLedgersAccount(ctx context.Context, i *dto.LedgersRequest) (
	*dto.ResponseDataList[[]dto.LedgerDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.ledgerUseCase.GetLedgersAccounts(req, *i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissioService.GetActions(req.Ctx, domain.LEDGER.ID)
	res.Body.Actions = actions
	return &res, nil
}

func (h *LedgerHandler) CreateLedger(ctx context.Context, i *dto.LedgerDataRequest) (
	*dto.ResponseData[dto.LedgerDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.ledgerUseCase.CreateLedger(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateNewAccount")
	}

	var response dto.ResponseData[dto.LedgerDto]
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateAccountLedgerSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	response.Body.Result = res
	return &response, nil
}
