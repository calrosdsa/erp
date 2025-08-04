package payment_rest

import (
	"context"
	// "erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	payment_ucase "erp/project/accounting/payment/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PaymentHandler struct {
	permissionService repository.PermissionService
	paymentUseCase    payment_ucase.PaymentUseCase
	sessionHelper     helpers.SessionHelper
	errorHelper       helpers.ErrorHelper
	locale            helpers.Locale
}

func NewPaymentHandler(
	api huma.API,
	middlewares huma.Middlewares,
	helpers *helpers.Helpers,
	paymentUseCase payment_ucase.PaymentUseCase,
	permissionService repository.PermissionService,
) {
	paths := NewPaymentPaths(domain.PAYMENT_BASE_ROUTE)
	tags := []string{"Payment"}
	h := PaymentHandler{
		paymentUseCase:    paymentUseCase,
		sessionHelper:     helpers.Session,
		errorHelper:       helpers.Error,
		locale:            helpers.Locale,
		permissionService: permissionService,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "create-payment",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Payment",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.CreatePayment)

	huma.Register(api, huma.Operation{
		OperationID:   "get-parties-type",
		Method:        http.MethodGet,
		Path:          paths.Parties,
		Summary:       "Get Parties Type",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPartiesType)

	huma.Register(api, huma.Operation{
		OperationID:   "get-payments",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Payments",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPayments)

	huma.Register(api, huma.Operation{
		OperationID:   "get payment detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Payment Detial",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPaymentDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "get payment actions",
		Method:        http.MethodGet,
		Path:          paths.AssociatedActions,
		Summary:       "Get Payment Actions",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPaymentAssociatedActions)

	huma.Register(api, huma.Operation{
		OperationID:   "update payment state",
		Method:        http.MethodPut,
		Path:          paths.UpdateState,
		Summary:       "Update Payment State",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdatePaymentState)

	huma.Register(api, huma.Operation{
		OperationID:   "payment-accouts",
		Method:        http.MethodGet,
		Path:          paths.PaymentAccounts,
		Summary:       "Payment Accounts",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPaymentAccounts)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-payment",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Payment",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditPayment)

	huma.Register(api,huma.Operation{
		OperationID: "export-payment",
		Method: http.MethodPost,
		Path: paths.Document,
		Summary: "Export Payment",
		Tags: tags,
		DefaultStatus: http.StatusOK,
		Middlewares: middlewares,
	},h.ExportPayment)
}

func (h *PaymentHandler) ExportPayment(ctx context.Context, i *dto.ExportDocumentRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes,err := h.paymentUseCase.ExportPayment(req,i.Body)
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

func (h *PaymentHandler) EditPayment(ctx context.Context, d *dto.EditPaymentRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.paymentUseCase.EditPayment(req, d.Body)
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

func (h *PaymentHandler) GetPaymentAccounts(ctx context.Context, i *struct{ dto.AuthParams }) (
	*dto.ResponseData[dto.PaymentAccountsDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentUseCase.GetPaymentAccounts(req)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissionService.GetActions(req.Ctx, domain.PAYMENT.ID)
	var response dto.ResponseData[dto.PaymentAccountsDto]
	response.Body.Actions = actions
	response.Body.Result = res
	return &response, nil
}

func (h *PaymentHandler) UpdatePaymentState(ctx context.Context, i *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.paymentUseCase.UpdatePaymentState(req, i)
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

func (h *PaymentHandler) GetPaymentAssociatedActions(ctx context.Context, i *struct{}) (
	*dto.ResponseData[any], error,
) {
	_, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	var res dto.ResponseData[any]
	res.Body.AssociatedActions = h.getExtraActions(ctx)
	return &res, nil
}

func (h *PaymentHandler) GetPayments(ctx context.Context, i *dto.RequestPayments) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.PaymentDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentUseCase.GetPayments(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissionService.GetActions(req.Ctx, domain.PAYMENT.ID)
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.PaymentDto]]
	response.Body.Actions = actions
	response.Body.PaginationResult = res
	return &response, nil
}

func (h *PaymentHandler) GetPaymentDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.PaymentDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentUseCase.GetPaymentDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permissionService.GetActions(req.Ctx, domain.PAYMENT.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.PaymentDetailDto]]
	response.Body.Actions = actions
	response.Body.Result = res
	response.Body.AssociatedActions = h.getExtraActions(ctx)
	return &response, nil
}

func (h *PaymentHandler) GetPartiesType(ctx context.Context, i *struct{ dto.AuthParams }) (
	*dto.ResponseData[[]dto.PartyTypeDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// req := &common.RequestContext{
	// 	LanguageCode:"en",
	// }
	res := h.paymentUseCase.GetAllowedParties(req)
	var response dto.ResponseData[[]dto.PartyTypeDto]
	response.Body.Result = res
	return &response, nil
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, i *dto.CreatePaymentRequest) (
	*dto.ResponseData[dto.PaymentDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.paymentUseCase.CreatePayment(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreatePayment")
	}
	var response dto.ResponseData[dto.PaymentDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatePaymentSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *PaymentHandler) getExtraActions(ctx context.Context) map[int][]dto.ActionDto {
	r := make(map[int][]dto.ActionDto)
	r[int(domain.LEDGER.ID)] = h.permissionService.GetActions(ctx, domain.LEDGER.ID)
	return r
}
