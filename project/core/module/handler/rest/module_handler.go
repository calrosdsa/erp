package module_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	module_ucase "erp/project/core/module/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ModuleHandler struct {
	sessionHelper       helpers.SessionHelper
	locale              helpers.Locale
	errorHelper         helpers.ErrorHelper
	moduleUcase module_ucase.ModuleUsecase
	permission          repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	moduleUcase module_ucase.ModuleUsecase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.MODULE_BASE_ROUTE
	tags := []string{"Module Handler"}
	path := NewPaths(base)
	h := ModuleHandler{
		sessionHelper:       helpers.Session,
		locale:              helpers.Locale,
		errorHelper:         helpers.Error,
		moduleUcase: moduleUcase,
		permission:          permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "search-entities",
		Method:        http.MethodGet,
		Summary:       "Search Entities",
		Tags:          tags,
		Path:          path.SearchEntities,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetEntitiesSearch)

	huma.Register(api, huma.Operation{
		OperationID:   "modules",
		Method:        http.MethodGet,
		Summary:       "Modules",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetModules)

	huma.Register(api, huma.Operation{
		OperationID:   "module",
		Method:        http.MethodGet,
		Summary:       "Module",
		Tags:          tags,
		Path:          path.Module,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetModule)

	huma.Register(api, huma.Operation{
		OperationID:   "module-detail",
		Method:        http.MethodGet,
		Summary:       "Module Detail",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetModuleDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "create-module",
		Method:        http.MethodPost,
		Summary:       "Create Module",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateModule)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-module",
		Method:        http.MethodPut,
		Path:          path.Base,
		Summary:       "Edit module",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditModule)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-module",
		Method:        http.MethodPut,
		Path:          path.UpdateStatus,
		Summary:       "Update Status module",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)
}

func (h *ModuleHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.moduleUcase.UpdateStatus(req, d)
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


func (h *ModuleHandler) EditModule(ctx context.Context, d *dto.ModuleDataRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.moduleUcase.EditModule(req, d.Body)
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

func (h *ModuleHandler) CreateModule(ctx context.Context, d *dto.ModuleDataRequest) (
	*dto.ResponseData[dto.ModuleDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.moduleUcase.CreateModule(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.ModuleDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}
func (h *ModuleHandler) GetModuleDetail(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ModuleDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.moduleUcase.GetModuleDetail(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.ModuleDetailDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.MODULE.ID)
	return &response, nil
}

func (h *ModuleHandler) GetModule(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ModuleDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.moduleUcase.GetModule(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.ModuleDetailDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.MODULE.ID)
	return &response, nil
}

func (h *ModuleHandler) GetModules(ctx context.Context, d *dto.ModulesRequest) (
	*dto.ResponseData[[]dto.ModuleDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.moduleUcase.GetModules(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.ModuleDto]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.MODULE.ID)
	return &response, nil
}


func (h *ModuleHandler) GetEntitiesSearch(ctx context.Context, d *dto.ModuleSearchRequest) (
	*dto.ResponseData[[]dto.EntityDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.moduleUcase.GetEntitiesSearch(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.EntityDto]
	response.Body.Result = res
	// response.Body.Actions = h.permission.GetActions(req.Ctx, domain.Mod.ID)
	return &response, nil
}
