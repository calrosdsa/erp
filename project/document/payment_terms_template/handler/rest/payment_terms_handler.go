package payment_terms_t_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	payment_terms_t_ucase "erp/project/document/payment_terms_template/usecase"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type paymentTermsTemplateHandler struct {
	sessionHelper             helpers.SessionHelper
	locale                    helpers.Locale
	errorHelper               helpers.ErrorHelper
	paymentTermsTemplateUcase payment_terms_t_ucase.PaymentTermsTemplateUcase
	permission                repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	paymentTermsTemplateUcase payment_terms_t_ucase.PaymentTermsTemplateUcase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.PAYMENT_TERMS_TEMPLATE_ROUTE
	tags := []string{"Payment Terms Template"}
	path := NewPaths(base)
	h := paymentTermsTemplateHandler{
		sessionHelper:             helpers.Session,
		locale:                    helpers.Locale,
		errorHelper:               helpers.Error,
		paymentTermsTemplateUcase: paymentTermsTemplateUcase,
		permission:                permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "payment-terms-template",
		Method:        http.MethodGet,
		Summary:       "Payment Terms Template",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPaymentTermsTemplates)
	huma.Register(api, huma.Operation{
		OperationID:   "payment-terms-template-details",
		Method:        http.MethodGet,
		Summary:       "Payment Terms Template Detail",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPaymentTermsTemplateDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "create-payment-terms-template",
		Method:        http.MethodPost,
		Summary:       "Create Payment Terms Template",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreatePaymentTermsTemplate)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-payment-terms-template",
		Method:        http.MethodPut,
		Path:          path.Base,
		Summary:       "Edit Payment Terms Template",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditPaymentTermsTemplate)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-payment-terms-template",
		Method:        http.MethodPut,
		Path:          path.UpdateStatus,
		Summary:       "Update Status Payment Terms-template",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "greet",
		Method:        http.MethodGet,
		Path:          path.Greet,
		Summary:       "Greet",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
	}, h.Greet)
}

func (h *paymentTermsTemplateHandler) Greet(ctx context.Context, d *struct {
	Name string `query:"name"`
}) (
	*struct {
		Body []byte
	}, error) {
	response := []byte(fmt.Sprintf("Hello, %s",d.Name))
	return &struct{Body []byte}{
		Body: response,
	}, nil
}

func (h *paymentTermsTemplateHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.paymentTermsTemplateUcase.UpdateStatus(req, *d)
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

func (h *paymentTermsTemplateHandler) EditPaymentTermsTemplate(ctx context.Context, d *dto.PaymentTermsTemplateDataRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.paymentTermsTemplateUcase.EditPaymentTermsTemplate(req, d.Body)
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

func (h *paymentTermsTemplateHandler) CreatePaymentTermsTemplate(ctx context.Context, d *dto.PaymentTermsTemplateDataRequest) (
	*dto.ResponseData[dto.PaymentTermsTemplateDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentTermsTemplateUcase.CreatePaymentTermsTemplate(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.PaymentTermsTemplateDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *paymentTermsTemplateHandler) GetPaymentTermsTemplateDetail(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.PaymentTermsTemplateDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentTermsTemplateUcase.GetPaymentTermsTemplateDetail(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.PaymentTermsTemplateDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.PAYMENT_TERMS_TEMPLATE.ID)
	return &response, nil
}

func (h *paymentTermsTemplateHandler) GetPaymentTermsTemplates(ctx context.Context, d *dto.PaymentTermsTemplateRequest) (
	*dto.ResponseDataList[[]dto.PaymentTermsTemplateDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentTermsTemplateUcase.GetPaymentTermsTemplates(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	res.Body.Actions = h.permission.GetActions(req.Ctx, domain.PAYMENT_TERMS_TEMPLATE.ID)
	return &res, nil
}
