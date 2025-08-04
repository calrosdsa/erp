package payment_terms_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	payment_terms_ucase "erp/project/document/payment_term/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type paymentTermsHandler struct {
	sessionHelper     helpers.SessionHelper
	locale            helpers.Locale
	errorHelper       helpers.ErrorHelper
	paymentTermsUcase payment_terms_ucase.PaymentTermsUcase
	permission        repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	paymentTermsUcase payment_terms_ucase.PaymentTermsUcase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.PAYMENT_TERMS_ROUTE
	tags := []string{"Payment Terms"}
	path := NewPaths(base)
	h := paymentTermsHandler{
		sessionHelper:     helpers.Session,
		locale:            helpers.Locale,
		errorHelper:       helpers.Error,
		paymentTermsUcase: paymentTermsUcase,
		permission:        permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "payment-terms",
		Method:        http.MethodGet,
		Summary:       "Payment Terms",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPaymentTerms)
	huma.Register(api, huma.Operation{
		OperationID:   "payment-terms-details",
		Method:        http.MethodGet,
		Summary:       "Payment Terms Detail",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPaymentTermsDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "create-payment-terms",
		Method:        http.MethodPost,
		Summary:       "Create Payment Terms",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreatePaymentTerms)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-payment-terms",
		Method:        http.MethodPut,
		Path:          path.Base,
		Summary:       "Edit Payment Terms",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditPaymentTerms)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-payment-terms",
		Method:        http.MethodPut,
		Path:          path.UpdateStatus,
		Summary:       "Update Status Payment Terms",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "payment-term-lines",
		Method:        http.MethodGet,
		Path:          path.Lines,
		Summary:       "Payment Term Lines",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPaymentTermLines)
}

func (h *paymentTermsHandler) GetPaymentTermLines(ctx context.Context, d *dto.RequestEntity) (
	*dto.ResponseData[[]dto.PaymentTermsLineDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res,err := h.paymentTermsUcase.GetPaymentTermLines(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdate")
	}
	var response dto.ResponseData[[]dto.PaymentTermsLineDto]
	response.Body.Result = res
	return &response, nil
}

func (h *paymentTermsHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.paymentTermsUcase.UpdateStatus(req, d)
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

func (h *paymentTermsHandler) EditPaymentTerms(ctx context.Context, d *dto.PaymentTermsDataRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.paymentTermsUcase.EditPaymentTerms(req, d.Body)
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

func (h *paymentTermsHandler) CreatePaymentTerms(ctx context.Context, d *dto.PaymentTermsDataRequest) (
	*dto.ResponseData[dto.PaymentTermsDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentTermsUcase.CreatePaymentTerms(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.PaymentTermsDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *paymentTermsHandler) GetPaymentTermsDetail(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.PaymentTermsDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentTermsUcase.GetPaymentTermsDetail(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.PaymentTermsDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.PAYMENT_TERMS.ID)
	return &response, nil
}

func (h *paymentTermsHandler) GetPaymentTerms(ctx context.Context, d *dto.PaymentTermsRequest) (
	*dto.ResponseDataList[[]dto.PaymentTermsDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentTermsUcase.GetPaymentTerms(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	res.Body.Actions = h.permission.GetActions(req.Ctx, domain.PAYMENT_TERMS.ID)
	return &res, nil
}
