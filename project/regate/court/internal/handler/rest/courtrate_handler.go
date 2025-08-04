package court_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	court_ucase "erp/project/regate/court/internal/usecase"
	regate_domain "erp/project/regate/internal/domain"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type CourtRateHandler struct {
	sessionHelper    helpers.SessionHelper
	locale           helpers.Locale
	errorHelper      helpers.ErrorHelper
	courtRateUseCase court_ucase.CourtRateUseCase
	permission       repository.PermissionService
}

func NewCourtRateHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	courtRateUseCase court_ucase.CourtRateUseCase,
) {
	paths := NewCourtRatePaths(domain.COURT_RATE_BASE_ROUTE)
	tag := []string{"Court Rate"}
	handler := CourtRateHandler{
		sessionHelper:    helpers.Session,
		locale:           helpers.Locale,
		errorHelper:      helpers.Error,
		courtRateUseCase: courtRateUseCase,
		permission:       permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "update-court-rates",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Update Court Rates",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateCourtRatesSchedule)

	huma.Register(api, huma.Operation{
		OperationID:   "get-court-rates",
		Method:        http.MethodGet,
		Path:          paths.Court,
		Summary:       "Get Court Rates",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetCourtRates)
	// huma.Register(api, huma.Operation{
	// 	OperationID:   "get court",
	// 	Method:        http.MethodGet,
	// 	Path:          paths.Detail,
	// 	Summary:       "Create User",
	// 	Tags:          tag,
	// 	DefaultStatus: http.StatusCreated,
	// 	Middlewares:   middlewares,
	// }, handler.GetCourt)
}

// func (h *CourtRateHandler) GetCourt(ctx context.Context, i *dto.RequestEntity) (
// 	*dto.EntityResponse[dto.ResultEntity[dto.CourtDto]], error,
// ) {
// 	req, err := h.sessionHelper.GetSession(ctx)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Not Authorized", err)
// 	}
// 	// h.sessionHelper.AppendPaginationParams(req, i)
// 	res, err := h.courtRateUseCase.GetCourt(req, i)
// 	if err != nil {
// 		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
// 	}
// 	actions := h.permission.GetActions(req.Ctx, regate_domain.COURT.ID)

// 	response := &dto.EntityResponse[dto.ResultEntity[dto.CourtDto]]{}
// 	response.Body.Result = res
// 	response.Body.Actions = actions
// 	return response, err
// }

func (h *CourtRateHandler) GetCourtRates(ctx context.Context, i *dto.RequestEntity) (
	*dto.ResponseData[[]dto.CourtRateDto], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.courtRateUseCase.GetCourtRates(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, regate_domain.COURT.ID)

	var response dto.ResponseData[[]dto.CourtRateDto]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, err
}

func (h *CourtRateHandler) UpdateCourtRatesSchedule(ctx context.Context, i *dto.UpdateCourtRatesRequest) (
	*dto.ResponseData[dto.CourtDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.courtRateUseCase.UpdateCourtRatesSchedule(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdate")
	}
	var response dto.ResponseData[dto.CourtDto]
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}
