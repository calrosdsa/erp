package rest_serial_no

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	serial_no_ucase "erp/project/stock/serial_no/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type SerialNoHandler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	serialUcase   serial_no_ucase.SerialNoUseCase
	permission    repository.PermissionService
}

func NewSerialHandler(
	api huma.API,
	helpers *helpers.Helpers,
	serialUcase serial_no_ucase.SerialNoUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.SERIAL_NO_BASE_ROUTE
	tags := []string{"Serian No"}
	path := NewSerialNoPaths(base)
	h := SerialNoHandler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		serialUcase:   serialUcase,
		permission:    permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "serial-nos",
		Method:        http.MethodGet,
		Summary:       "Serial Nos",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetSerialNos)
	huma.Register(api, huma.Operation{
		OperationID:   "serial-no",
		Method:        http.MethodGet,
		Summary:       "Serial No",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetSerialNo)

	huma.Register(api, huma.Operation{
		OperationID:   "serial-no-transactions",
		Method:        http.MethodGet,
		Summary:       "Serial No Transactions",
		Tags:          tags,
		Path:          path.SerialNoTransactions,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetSerialNoTransactions)
}
func (h *SerialNoHandler) GetSerialNoTransactions(ctx context.Context, d *dto.RequestSerialNoTransactions) (
	*dto.ResponseData[[]dto.SerialNoTransactionDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.serialUcase.GetSerialNoTransactions(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.SerialNoTransactionDto]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.SERIAL_NO.ID)
	return &response, nil
}

func (h *SerialNoHandler) GetSerialNo(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.SerialNoDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.serialUcase.GetSerialNo(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.SerialNoDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.SERIAL_NO.ID)
	return &response, nil
}

func (h *SerialNoHandler) GetSerialNos(ctx context.Context, d *dto.RequestSerialNos) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.SerialNoDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.serialUcase.GetSerialNos(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.SerialNoDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.SERIAL_NO.ID)
	return &response, nil
}
