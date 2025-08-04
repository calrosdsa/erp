package cash_outflow_rest

import (
	context "context"
	dto "erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	cash_outflow_ucase "erp/project/accounting/cash_outflow/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	usecase       cash_outflow_ucase.CashOutflowUseCase
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	usecase cash_outflow_ucase.CashOutflowUseCase,
) {
	base := domain.CASH_OUTFLOW_ROUTE
	tags := []string{"Cash Outflow"}
	h := handler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		usecase:       usecase,
		permission:    permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "cash-outflow",
		Method:        http.MethodGet,
		Summary:       "Cash Outflow",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetList)
	huma.Register(api, huma.Operation{
		OperationID:   "cash-outflow-detail",
		Method:        http.MethodGet,
		Summary:       "Cash Outflow Detail",
		Tags:          tags,
		Path:          base + "/detail/{id}",
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.Get)

	huma.Register(api, huma.Operation{
		OperationID:   "create-cash-outflow",
		Method:        http.MethodPost,
		Summary:       "Create Cash Outflow",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.Create)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-cash-outflow",
		Method:        http.MethodPut,
		Path:          base,
		Summary:       "Edit Cash Outflow",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.Edit)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-cash-outflow",
		Method:        http.MethodPut,
		Path:          base + "/update-status",
		Summary:       "Update Status Cash Outflow",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)

	huma.Register(api,huma.Operation{
		OperationID: "export-cash-outflow",
		Method: http.MethodPost,
		Path:  base + "/export/document",
		Summary: "Export Cash Outflow",
		Tags: tags,
		DefaultStatus: http.StatusOK,
		Middlewares: middlewares,
	},h.ExportCashOutflow)
}

func (h *handler) ExportCashOutflow(ctx context.Context, i *dto.ExportDocumentRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes,err := h.usecase.ExportCashOutflow(req,i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			writer := ctx.BodyWriter()
			writer.Write(bytes)
		},
	}, nil
}

// Create provides a function with given fields: ctx, d

func (m *handler) Create(ctx context.Context, d *dto.CashOutflowDataRequest) (_a0 *dto.ResponseData[dto.CashOutflowDto], _a1 error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.Create(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.CashOutflowDto]
	response.Body.Result = res
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (m *handler) Edit(ctx context.Context, d *dto.CashOutflowDataRequest) (_a0 *dto.ResponseMessage, _a1 error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.EditCashOutflow(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

// Get provides a function with given fields: ctx, d

func (m *handler) Get(ctx context.Context, d *dto.RequestEntity) (_a0 *dto.EntityResponse[dto.ResultEntity[dto.CashOutflowDto]], _a1 error) {

	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetCashOutflow(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.CashOutflowDto]]
	response.Body.Result = res
	response.Body.Actions = m.permission.GetActions(req.Ctx, domain.CASH_OUTFLOW.ID)
	response.Body.AssociatedActions = m.getExtraActions(ctx)
	return &response, nil

}

// GetList provides a function with given fields: ctx, d

func (m *handler) GetList(ctx context.Context, d *dto.CashOutflowsRequest) (_a0 *dto.ResponseDataList[[]dto.CashOutflowDto], _a1 error) {

	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetCashOutflows(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	res.Body.Actions = m.permission.GetActions(req.Ctx, domain.CASH_OUTFLOW.ID)
	return &res, nil

}

// UpdateStatus provides a function with given fields: ctx, d

func (m *handler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (_a0 *dto.ResponseMessage, _a1 error) {

	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.UpdateStatus(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdate")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (h *handler) getExtraActions(ctx context.Context) map[int][]dto.ActionDto {
	var ids []int64
	ids = append(ids, domain.COST_CENTER.ID, domain.PROJECT.ID, domain.LEDGER.ID,
		domain.GENERAL_LEDGER.ID)
	r := h.permission.GetEntitiesActions(ctx, ids)
	return r
}
