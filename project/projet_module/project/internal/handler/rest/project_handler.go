package rest_project

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	project_ucase "erp/project/projet_module/project/internal/usecase"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ProjectHandler struct {
	sessionHelper   helpers.SessionHelper
	locale          helpers.Locale
	errorHelper     helpers.ErrorHelper
	projectUseCase project_ucase.ProjectUseCase
	permission repository.PermissionService
}

func NewProjectHandler(
	api huma.API,
	helpers *helpers.Helpers,
	projectUseCase project_ucase.ProjectUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.PROJECT_BASE_ROUTE
	tags := []string{"Project"}
	path := NewProjectPaths(base)
	h := ProjectHandler{
		sessionHelper:   helpers.Session,
		locale:          helpers.Locale,
		errorHelper:     helpers.Error,
		projectUseCase: projectUseCase,
		permission: permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "projects",
		Method:        http.MethodGet,
		Summary:       "Projects",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetProjects)
	huma.Register(api, huma.Operation{
		OperationID:   "project",
		Method:        http.MethodGet,
		Summary:       "Project",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetProject)

	huma.Register(api, huma.Operation{
		OperationID:   "create-cost-project",
		Method:        http.MethodPost,
		Summary:       "Create Project",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateProject)

	huma.Register(api,huma.Operation{
		OperationID:   "test-request",
		Method:        http.MethodDelete,
		Summary:       "Test Request",
		Tags:          tags,
		Path:          path.TestRequest,
		DefaultStatus: http.StatusOK,
	},h.ReceiveRequest)
}

// https://fortixdrconnectlat2.console.ensilo.com/management-rest/inventory/delete-collectors
// ?states=Disconnected&collectorGroups=NOMBRE_DEL_GRUPO
func (h *ProjectHandler) ReceiveRequest(ctx context.Context,d *struct{
	State string `query:"states"`
	CollectorGroups string `query:"collectorGroups"`
})(*dto.ResponseMessage,error){
	fmt.Println("STATE",d.State)
	fmt.Println("CollectorGroups",d.CollectorGroups)
	var response dto.ResponseMessage
	response.Body.Message  = "Success"
	return &response,nil
}
func (h *ProjectHandler) CreateProject(ctx context.Context, d *dto.CreateProjectRequest) (
	*dto.ResponseData[dto.ProjectDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.projectUseCase.CreateProject(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.ProjectDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ProjectHandler) GetProject(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ProjectDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.projectUseCase.GetProject(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.ProjectDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(ctx,domain.PROJECT.ID)
	return &response, nil
}

func (h *ProjectHandler) GetProjects(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.ProjectDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.projectUseCase.GetProjects(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.ProjectDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(ctx,domain.PROJECT.ID)
	return &response, nil
}
