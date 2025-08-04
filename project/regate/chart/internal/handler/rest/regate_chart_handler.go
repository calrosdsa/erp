package chart_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	chart_ucase "erp/project/regate/chart/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ChartHandler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	chartUseCase  chart_ucase.ChartUseCase
}

func NewChartHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	chartUseCase chart_ucase.ChartUseCase,
) {
	paths := NewChartPaths(domain.CHART_BASE_ROUTE)
	tag := []string{"Chart"}
	handler := ChartHandler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		chartUseCase:  chartUseCase,
		permission:    permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "get-chart-data",
		Method:        http.MethodPost,
		Path:          paths.Chart,
		Summary:       "Get Chart Data",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetChartData)

	huma.Register(api, huma.Operation{
		OperationID:   "get-chart-dashboard-data",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Get Chart Dashboard Data",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetDashboardData)
}

func (h *ChartHandler) GetDashboardData(ctx context.Context, i *dto.ChartDashboardDataRequest) (
	*dto.ResponseData[dto.ChartDashboardData], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.chartUseCase.GetDashboardData(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	// actions := h.permission.GetActions(req.Ctx, regate_domain.EVENT_BOOKING.ID)

	var response dto.ResponseData[dto.ChartDashboardData]
	response.Body.Result = res
	// response.Body.Actions = actions
	return &response, err
}

func (h *ChartHandler) GetChartData(ctx context.Context, i *dto.ChartDataRequest) (
	*dto.ResponseData[[]dto.ChartDataDto], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.chartUseCase.GetChartData(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	// actions := h.permission.GetActions(req.Ctx, regate_domain.EVENT_BOOKING.ID)

	var response dto.ResponseData[[]dto.ChartDataDto]
	response.Body.Result = res
	// response.Body.Actions = actions
	return &response, err
}
