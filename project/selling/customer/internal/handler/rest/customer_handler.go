package customer_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	customer_ucase "erp/project/selling/customer/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type CustomerHandler struct {
	sessionHelper helpers.SessionHelper
	errorHelper   helpers.ErrorHelper
	custumerUcase customer_ucase.CustomerUseCase
	permission    repository.PermissionService
	locale        helpers.Locale
}

func NewCustomerHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	customerUcase customer_ucase.CustomerUseCase,
) {
	paths := NewCustomerPaths(domain.CUSTOMER_BASE_ROUTE)
	tags := []string{"Customer"}
	h := CustomerHandler{
		sessionHelper: helpers.Session,
		custumerUcase: customerUcase,
		permission:    permission,
		errorHelper:   helpers.Error,
		locale:        helpers.Locale,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "get customer",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Retrieve customer",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCustomer)
	huma.Register(api, huma.Operation{
		OperationID:   "get customers",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Retrieve customers",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCustomers)

	huma.Register(api, huma.Operation{
		OperationID:   "create customer",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create customer",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateCustomer)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-customer",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit customer",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditCustomer)

	huma.Register(api, huma.Operation{
		OperationID:   "customer types",
		Method:        http.MethodGet,
		Path:          paths.CustomerTypes,
		Summary:       "Customer types",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCustomerTypes)
	huma.Register(api, huma.Operation{
		OperationID:   "update-status-customer",
		Method:        http.MethodPut,
		Path:          paths.UpdateStatus,
		Summary:       "Update Status Customer",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)
}
func (h *CustomerHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.custumerUcase.UpdateStatus(req, d)
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

func (h *CustomerHandler) EditCustomer(ctx context.Context, d *dto.CustomerDataRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.custumerUcase.EditCustomer(req, d.Body)
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

func (h *CustomerHandler) GetCustomerTypes(ctx context.Context, i *struct {
	dto.AuthParams
}) (
	*dto.EntityResponse[dto.ResultEntity[[]dto.CustomerType]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res := h.custumerUcase.GetCustomerTypes(req)

	var response dto.EntityResponse[dto.ResultEntity[[]dto.CustomerType]]
	response.Body.Result.Entity = res
	return &response, nil
}

func (h *CustomerHandler) GetCustomer(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.CustomerDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.custumerUcase.GetCustomerDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.CUSTOMER.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.CustomerDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *CustomerHandler) CreateCustomer(ctx context.Context, i *dto.CustomerDataRequest) (
	*dto.ResponseData[dto.CustomerDto], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res,err := h.custumerUcase.CreateCustomer(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateCustomer")
	}
	var response dto.ResponseData[dto.CustomerDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateCustomerSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *CustomerHandler) GetCustomers(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CustomerDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.custumerUcase.GetCustomers(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.CUSTOMER.ID)

	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.CustomerDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}
