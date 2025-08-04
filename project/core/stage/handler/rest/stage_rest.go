package stage_rest

import (
	context "context"
	dto "erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	stage_ucase "erp/project/core/stage/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	usecase       stage_ucase.StageUseCase
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	usecase stage_ucase.StageUseCase,
) {
	base := domain.STAGE_BASE_ROUTE
	tags := []string{"Stage"}
	h := handler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		usecase:       usecase,
		permission:    permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "stage",
		Method:        http.MethodGet,
		Summary:       "Stage",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetStages)

	huma.Register(api, huma.Operation{
		OperationID:   "create-stage",
		Method:        http.MethodPost,
		Summary:       "Create Stage",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateStage)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-stage",
		Method:        http.MethodPut,
		Path:          base,
		Summary:       "Edit Stage",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditStage)

	huma.Register(api, huma.Operation{
		OperationID:   "stage-transition",
		Method:        http.MethodPut,
		Path:          base + "/transition",
		Summary:       "Stage Transition",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.StageTransition)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-stage",
		Method:        http.MethodDelete,
		Path:          base,
		Summary:       "Delete Stage",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.DeleteStage)
}


func (m *handler) DeleteStage(ctx context.Context, d *dto.DeleteRequest) (*dto.ResponseMessage, error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.DeleteStage(req, d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToDelete")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.DeletedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (m *handler) StageTransition(ctx context.Context, d *dto.StageTransitionRequest) (*dto.ResponseMessage, error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.StageTransition(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (m *handler) CreateStage(ctx context.Context, d *dto.StageDataRequest) (*dto.ResponseData[dto.StageDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.CreateStage(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.StageDto]
	response.Body.Result = res
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (m *handler) EditStage(ctx context.Context, d *dto.StageDataRequest) (*dto.ResponseMessage, error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.EditStage(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (m *handler) GetStages(ctx context.Context, d *dto.StagesRequest) (*dto.ResponseDataList[[]dto.StageDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetStages(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	res.Body.Actions = m.permission.GetActions(ctx,domain.STAGE.ID)
	// res.Body.AssociatedActions = m.getActions(ctx)
	return &res, nil

}

func (h *handler) getActions(ctx context.Context) map[int][]dto.ActionDto {
	var ids []int64
	ids = append(ids, domain.STAGE.ID)
	r := h.permission.GetEntitiesActions(ctx, ids)
	return r
}
