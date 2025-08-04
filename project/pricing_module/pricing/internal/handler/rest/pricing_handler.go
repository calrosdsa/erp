package pricing_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	pricing_ucase "erp/project/pricing_module/pricing/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PricingHandler struct {
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	pricingUseCase pricing_ucase.PricingUseCase
	pricingGeneratorUseCase pricing_ucase.PricingGeneratorUcase
	permission     repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	pricingUseCase pricing_ucase.PricingUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	pricingGeneratorUseCase pricing_ucase.PricingGeneratorUcase,
) {
	base := domain.PRICING_BASE_ROUTE
	tags := []string{"Pricing"}
	path := NewPaths(base)
	h := PricingHandler{
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
		pricingUseCase: pricingUseCase,
		permission:     permission,
		pricingGeneratorUseCase: pricingGeneratorUseCase,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "pricings",
		Method:        http.MethodGet,
		Summary:       "Pricings",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPricings)
	huma.Register(api, huma.Operation{
		OperationID:   "pricing",
		Method:        http.MethodGet,
		Summary:       "Pricing",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPricing)

	huma.Register(api, huma.Operation{
		OperationID:   "create-pricing",
		Method:        http.MethodPost,
		Summary:       "Create pricing",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreatePricing)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-pricing",
		Method:        http.MethodPut,
		Path:          path.Base,
		Summary:       "Edit Pricing",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditPricing)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-pricing",
		Method:        http.MethodPut,
		Path:          path.UpdateStatus,
		Summary:       "Update Status Pricing",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "pricing-generate-po",
		Method:        http.MethodPost,
		Path:          path.GeneratePo,
		Summary:       "Pricing Generate Po",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GeneratePo)

	huma.Register(api, huma.Operation{
		OperationID:   "pricing-generate-quotation",
		Method:        http.MethodPost,
		Path:          path.GenerateQuotation,
		Summary:       "Pricing Generate Quotation",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GenerateQuotation)
}

func (h *PricingHandler) GenerateQuotation(ctx context.Context, d *dto.PricingDataRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.pricingGeneratorUseCase.GenerateQuotation(req, d)
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


func (h *PricingHandler) GeneratePo(ctx context.Context, d *dto.PricingDataRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.pricingGeneratorUseCase.GeneratePo(req, d)
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


func (h *PricingHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.pricingUseCase.UpdateStatus(req, d)
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

func (h *PricingHandler) EditPricing(ctx context.Context, d *dto.EditPricingRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.pricingUseCase.EditPricing(req, d)
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

func (h *PricingHandler) CreatePricing(ctx context.Context, d *dto.CreatePricingRequest) (
	*dto.ResponseData[dto.PricingDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.pricingUseCase.CreatePricing(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.PricingDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *PricingHandler) GetPricing(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.PricingDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.pricingUseCase.GetPricing(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.PricingDetailDto]]
	response.Body.Result = res
	// response.Body.Actions = h.permission.GetActions(req.Ctx, domain.PRICING.ID)
	response.Body.AssociatedActions = h.getExtraActions(ctx)
	return &response, nil
}

func (h *PricingHandler) GetPricings(ctx context.Context, d *dto.RequestPricings) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.PricingDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.pricingUseCase.GetPricings(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.PricingDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.PRICING.ID)
	return &response, nil
}

func (h *PricingHandler) getExtraActions(ctx context.Context) map[int][]dto.ActionDto {
	r := make(map[int][]dto.ActionDto)
	r[int(domain.PURCHASE_ORDER.ID)] = h.permission.GetActions(ctx, domain.PURCHASE_ORDER.ID)
	r[int(domain.QUOTATION.ID)] = h.permission.GetActions(ctx, domain.QUOTATION.ID)
	r[int(domain.PRICING.ID)] = h.permission.GetActions(ctx, domain.PRICING.ID)
	return r
}
