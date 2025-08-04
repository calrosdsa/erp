package supplier_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	supplier_ucase "erp/project/buying/supplier/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type SupplierHandler struct {
	sessionHelper   helpers.SessionHelper
	errorHelper     helpers.ErrorHelper
	locale          helpers.Locale
	permission      repository.PermissionService
	supplierUseCase supplier_ucase.SupplierUseCase
}

func NewSupplierHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	supplierUseCase supplier_ucase.SupplierUseCase,
	permission repository.PermissionService,
) {
	paths := NewSupplierPaths(domain.SUPPLIER_BASE_ROUTE)
	tags := []string{"Supplier"}
	h := SupplierHandler{
		sessionHelper:   helpers.Session,
		errorHelper:     helpers.Error,
		supplierUseCase: supplierUseCase,
		locale:          helpers.Locale,
		permission:      permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "get supplier",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Retrieve supplier",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetSupplier)

	huma.Register(api, huma.Operation{
		OperationID:   "get suppliers",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Retrieve suppliers",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetSuppliers)

	huma.Register(api, huma.Operation{
		OperationID:   "create supplier",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create supplier",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateSupplier)
	huma.Register(api, huma.Operation{
		OperationID:   "edit-supplier",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Supplier",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditSupplier)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-supplier",
		Method:        http.MethodPut,
		Path:          paths.UpdateStatus,
		Summary:       "Update Status Supplier",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)
}
func (h *SupplierHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.supplierUseCase.UpdateStatus(req, d)
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

func (h *SupplierHandler) EditSupplier(ctx context.Context, d *dto.SupplierDataRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.supplierUseCase.EditSupplier(req, d.Body)
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

func (h *SupplierHandler) GetSupplier(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.SupplierDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.supplierUseCase.GetSupplier(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.SUPPLIER.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.SupplierDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *SupplierHandler) CreateSupplier(ctx context.Context, i *dto.SupplierDataRequest) (
	*dto.ResponseData[dto.SupplierDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res,err := h.supplierUseCase.CreateSupplier(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseData[dto.SupplierDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *SupplierHandler) GetSuppliers(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.SupplierDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.supplierUseCase.GetSuppliers(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.SUPPLIER.ID)

	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.SupplierDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}
