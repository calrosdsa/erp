package rest_task

import (
	context "context"
	dto "erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	task_ucase "erp/project/projet_module/task/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	usecase       task_ucase.TaskUseCase
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	usecase task_ucase.TaskUseCase,
) {
	base := domain.TASK_BASE_ROUTE
	tags := []string{"Task"}
	h := handler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		usecase:       usecase,
		permission:    permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "task",
		Method:        http.MethodGet,
		Summary:       "Task",
		Description:   "Retrieve a paginated list of tasks with optional filtering and sorting capabilities. Returns task summary information including status, assignee, and key dates.",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetTasks)
	huma.Register(api, huma.Operation{
		OperationID:   "task-detail",
		Method:        http.MethodGet,
		Summary:       "Task Detail",
		Description:   "Retrieve comprehensive details for a specific task by ID. Returns full task information including project details, assignee information, and associated actions available to the current user.",
		Tags:          tags,
		Path:          base + "/detail/{id}",
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetTask)

	huma.Register(api, huma.Operation{
		OperationID:   "create-task",
		Method:        http.MethodPost,
		Summary:       "Create Task",
		Description:   "Create a new task record with the provided task data. Validates input data, assigns initial status, and returns the created task with success message.",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateTask)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-task",
		Method:        http.MethodPut,
		Path:          base,
		Summary:       "Edit Task",
		Description:   "Update an existing task record with new data. Validates the provided changes and updates the task information while maintaining audit trail.",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditTask)


	huma.Register(api, huma.Operation{
		OperationID:   "task-transition",
		Method:        http.MethodPut,
		Path:          base + "/transition",
		Summary:       "Task Transition",
		Description:   "Transition a task through its lifecycle states (e.g., from prospect to qualified, won, or lost). Validates state transition rules and updates task status accordingly.",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.TaskTransition)
}
func (m *handler) TaskTransition(ctx context.Context, d *dto.EntityTransitionRequest) (*dto.ResponseMessage, error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.TaskTransition(req, d.Body)
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


func (m *handler) CreateTask(ctx context.Context, d *dto.TaskDataRequest) (*dto.ResponseData[dto.TaskDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.CreateTask(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.TaskDto]
	response.Body.Result = res
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (m *handler) EditTask(ctx context.Context, d *dto.TaskDataRequest) (*dto.ResponseMessage, error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.EditTask(req, d.Body)
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

func (m *handler) GetTask(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.TaskDetailDto]], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetTask(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.TaskDetailDto]]
	response.Body.Result = res
	response.Body.Actions = m.permission.GetActions(ctx, domain.TASK.ID)
	response.Body.AssociatedActions = m.getActions(ctx)
	return &response, nil

}

func (m *handler) GetTasks(ctx context.Context, d *dto.TasksRequest) (*dto.ResponseDataList[[]dto.TaskDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetTasks(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	res.Body.Actions = m.permission.GetActions(ctx, domain.TASK.ID)
	return &res, nil

}

func (h *handler) getActions(ctx context.Context) map[int][]dto.ActionDto {
	var ids []int64
	ids = append(ids, domain.PROJECT.ID)
	ids = append(ids, domain.CONTACT.ID)

	r := h.permission.GetEntitiesActions(ctx, ids)

	return r
}
