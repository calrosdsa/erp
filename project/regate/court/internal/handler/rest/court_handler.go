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

type CourtHandler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	courtUseCase  court_ucase.CourtUseCase
	permission    repository.PermissionService
}

func NewCourtHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	courtUseCase court_ucase.CourtUseCase,
) {
	paths := NewCourtPaths(domain.COURT_BASE_ROUTE)
	tag := []string{"Court"}
	handler := CourtHandler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		courtUseCase:  courtUseCase,
		permission:    permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "create-court",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Court",
		Tags:          tag,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateCourt)

	huma.Register(api, huma.Operation{
		OperationID:   "courts",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "get courts",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetCourts)
	huma.Register(api, huma.Operation{
		OperationID:   "court",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Create User",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetCourt)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-court",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit court",
		Tags:          tag,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.EditCourt)
}
func (h *CourtHandler) EditCourt(ctx context.Context, d *dto.EditCourtRequest) (*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.courtUseCase.EditCourt(req, d.Body)
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

func (h *CourtHandler) GetCourt(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.CourtDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.courtUseCase.GetCourt(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, regate_domain.COURT.ID)

	response := &dto.EntityResponse[dto.ResultEntity[dto.CourtDto]]{}
	response.Body.Result = res
	response.Body.Actions = actions
	return response, err
}

func (h *CourtHandler) GetCourts(ctx context.Context, i *dto.CourtsRequest) (
	*dto.ResponseDataList[[]dto.CourtDto],error) {
		req, err := h.sessionHelper.GetSession(ctx)
		if err != nil {
			return nil, huma.Error400BadRequest("Not Authorized", err)
		}
		res, err := h.courtUseCase.GetCourts(req, *i)
		if err != nil {
			return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
		}
		res.Body.Actions = h.permission.GetActions(req.Ctx, domain.BANK.ID)
		return &res, nil
}

func (h *CourtHandler) CreateCourt(ctx context.Context, i *dto.CreateCourtRequest) (
	*dto.ResponseData[dto.CourtDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.courtUseCase.CreateCourt(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.CourtDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}
