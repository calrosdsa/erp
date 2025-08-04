package charges_template_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	charge_template_ucase "erp/project/accounting/charges_template/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ChargeTemplateHandler struct {
	sessionHelper       helpers.SessionHelper
	locale              helpers.Locale
	errorHelper         helpers.ErrorHelper
	chargeTemplateUcase charge_template_ucase.ChargesTemplateUseCase
	permission          repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	chargeTemplateUcase charge_template_ucase.ChargesTemplateUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.CHARGE_TEMPLATE_BASE_ROUTE
	tags := []string{"Charges Template"}
	path := NewPaths(base)
	h := ChargeTemplateHandler{
		sessionHelper:       helpers.Session,
		locale:              helpers.Locale,
		errorHelper:         helpers.Error,
		chargeTemplateUcase: chargeTemplateUcase,
		permission:          permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "charges-templates",
		Method:        http.MethodGet,
		Summary:       "Charges Templates",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetChargesTemplates)
	huma.Register(api, huma.Operation{
		OperationID:   "charges-template",
		Method:        http.MethodGet,
		Summary:       "Charges Template",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetChargesTemplate)

	huma.Register(api, huma.Operation{
		OperationID:   "create-charge-template",
		Method:        http.MethodPost,
		Summary:       "Create Charge Template",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateChargesTemplate)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-charges-template",
		Method:        http.MethodPut,
		Path:          path.Base,
		Summary:       "Edit Charges Template",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditChargesTemplate)
}
func (h *ChargeTemplateHandler) EditChargesTemplate(ctx context.Context, d *dto.EditChargesTemplateRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.chargeTemplateUcase.EditChargesTemplate(req, d)
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

func (h *ChargeTemplateHandler) CreateChargesTemplate(ctx context.Context, d *dto.CreateChargesTemplateRequest) (
	*dto.ResponseData[dto.ChargesTemplateDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.chargeTemplateUcase.CreateChargesTemplate(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.ChargesTemplateDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ChargeTemplateHandler) GetChargesTemplate(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ChargesTemplateDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.chargeTemplateUcase.GetChargesTemplate(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.ChargesTemplateDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.CHARGES_TEMPLATE.ID)
	return &response, nil
}

func (h *ChargeTemplateHandler) GetChargesTemplates(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.ChargesTemplateDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.chargeTemplateUcase.GetChargesTemplates(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.ChargesTemplateDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.CHARGES_TEMPLATE.ID)
	return &response, nil
}
