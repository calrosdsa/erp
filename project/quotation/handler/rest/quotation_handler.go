package rest_quotation

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	quotation_ucase "erp/project/quotation/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type quotationHandler struct {
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	quotationUcase quotation_ucase.QuotationUseCase
	permission     repository.PermissionService
}

func NewQuotationHandler(
	api huma.API,
	helpers *helpers.Helpers,
	quotationUcase quotation_ucase.QuotationUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.QUOTATION_BASE_ROUTE
	tags := []string{"Quotation"}
	path := NewQuotationPaths(base)
	h := quotationHandler{
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
		quotationUcase: quotationUcase,
		permission:     permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "quotations",
		Method:        http.MethodGet,
		Summary:       "Quotations",
		Tags:          tags,
		Path:          path.QuotationParty,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetQuotations)
	huma.Register(api, huma.Operation{
		OperationID:   "quotation",
		Method:        http.MethodGet,
		Summary:       "Quotation",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetQuotation)

	huma.Register(api, huma.Operation{
		OperationID:   "create-quotation",
		Method:        http.MethodPost,
		Summary:       "Create Quotation",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateQuotation)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-quotation",
		Method:        http.MethodPut,
		Summary:       "Edit Quotation",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditQuotation)

	huma.Register(api, huma.Operation{
		OperationID:   "update-quotation-status",
		Method:        http.MethodPut,
		Summary:       "Update Quotation Status",
		Path:          path.UpdateStatus,
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)
}
func (h *quotationHandler) UpdateStatus(ctx context.Context, i *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.quotationUcase.UpdateStatus(req, i)
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

func (h *quotationHandler) EditQuotation(ctx context.Context, d *dto.EditQuotationRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.quotationUcase.EditQuotation(req, d.Body)
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
func (h *quotationHandler) CreateQuotation(ctx context.Context, d *dto.CreateQuotationRequest) (
	*dto.ResponseData[dto.QuotationDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.quotationUcase.CreateQuotation(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseData[dto.QuotationDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *quotationHandler) GetQuotation(ctx context.Context, d *dto.RequestEntityWithParty) (
	*dto.EntityResponse[dto.ResultEntity[dto.QuotationDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	entity, err := h.quotationUcase.GetQuotationEntity(d.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	res, err := h.quotationUcase.GetQuotation(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.QuotationDetailDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.QUOTATION.ID)
	response.Body.AssociatedActions = h.getExtraActions(req.Ctx, entity)
	return &response, nil
}

func (h *quotationHandler) GetQuotations(ctx context.Context, d *dto.RequestQuotations) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.QuotationDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.quotationUcase.GetQuotations(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.QuotationDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.QUOTATION.ID)
	return &response, nil
}

func (h *quotationHandler) getExtraActions(ctx context.Context, entity domain.EntityTemplate) map[int][]dto.ActionDto {
	r := make(map[int][]dto.ActionDto)
	switch entity {
	case domain.SUPPLIER_QUOTATION:
		r[int(domain.PURCHASE_ORDER.ID)] = h.permission.GetActions(ctx, domain.PURCHASE_ORDER.ID)
		r[int(domain.QUOTATION.ID)] = h.permission.GetActions(ctx, domain.QUOTATION.ID)
	case domain.QUOTATION:
		r[int(domain.SALE_ORDER.ID)] = h.permission.GetActions(ctx, domain.SALE_ORDER.ID)
	}
	return r
}
