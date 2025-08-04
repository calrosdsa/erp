package rest_cost_center

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	cost_center_ucase "erp/project/accounting/cost_center/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type CostCenterHandler struct {
	sessionHelper   helpers.SessionHelper
	locale          helpers.Locale
	errorHelper     helpers.ErrorHelper
	costCenterUcase cost_center_ucase.CostCenterUseCase
	permission repository.PermissionService
}

func NewCostCenterHandler(
	api huma.API,
	helpers *helpers.Helpers,
	costCenterUcase cost_center_ucase.CostCenterUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.COST_CENTER_BASE_ROUTE
	tags := []string{"Cost Center"}
	path := NewCostCenterPaths(base)
	h := CostCenterHandler{
		sessionHelper:   helpers.Session,
		locale:          helpers.Locale,
		errorHelper:     helpers.Error,
		costCenterUcase: costCenterUcase,
		permission: permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "cost-centers",
		Method:        http.MethodGet,
		Summary:       "Cost Centers",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCostCenters)
	huma.Register(api, huma.Operation{
		OperationID:   "cost-center",
		Method:        http.MethodGet,
		Summary:       "Cost Center",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCostCenter)

	huma.Register(api, huma.Operation{
		OperationID:   "create-cost-center",
		Method:        http.MethodPost,
		Summary:       "Create Cost Center",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateCostCenter)
}
func (h *CostCenterHandler) CreateCostCenter(ctx context.Context, d *dto.CreateCostCenterRequet) (
	*dto.ResponseData[dto.CostCenterDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.costCenterUcase.CreateCostCenter(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseData[dto.CostCenterDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *CostCenterHandler) GetCostCenter(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.CostCenterDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.costCenterUcase.GetCostCenter(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.CostCenterDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx,domain.COST_CENTER.ID)
	return &response, nil
}

func (h *CostCenterHandler) GetCostCenters(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CostCenterDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.costCenterUcase.GetCostCenters(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.CostCenterDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(req.Ctx,domain.COST_CENTER.ID)
	return &response, nil
}
