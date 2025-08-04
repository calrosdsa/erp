package role_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	role_ucase "erp/project/manage/role/internal/usecase"

	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type RoleHandler struct {
	jwtHelper     helpers.JwtHelper
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	roleUseCase   role_ucase.RoleUseCase
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
}

func NewRoleHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	roleUseCase role_ucase.RoleUseCase,
) {
	paths := NewRolePaths(domain.ROLE_BASE_ROUTE)
	tag := "Role"
	handler := RoleHandler{
		jwtHelper:     helpers.Jwt,
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		roleUseCase:   roleUseCase,
		permission:    permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "get roles",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get roles",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetRoles)
	huma.Register(api, huma.Operation{
		OperationID:   "get role definitions",
		Method:        http.MethodGet,
		Path:          paths.RoleDefinitions,
		Summary:       "Get role definitions",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetRoleDefinitions)

	huma.Register(api, huma.Operation{
		OperationID:   "get entity actions",
		Method:        http.MethodGet,
		Path:          paths.EntityActions,
		Summary:       "Get entity actions",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetEntityActions)

	huma.Register(api, huma.Operation{
		OperationID:   "get role",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get role",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetRole)

	huma.Register(api, huma.Operation{
		OperationID:   "update role permision action",
		Method:        http.MethodPost,
		Path:          paths.PermissionActions,
		Summary:       "Update Permission Actions",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateRolePermissionAction)

	huma.Register(api, huma.Operation{
		OperationID:   "create-role",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Role",
		Tags:          []string{tag},
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateRole)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-role",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Role",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditRole)
}

func (m *RoleHandler) EditRole(ctx context.Context, d *dto.RoleRequestData) (_a0 *dto.ResponseMessage, _a1 error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.roleUseCase.EditRole(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}


func (h *RoleHandler) CreateRole(ctx context.Context, i *dto.RoleRequestData) (*dto.ResponseData[dto.RoleDto], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.roleUseCase.CreateRole(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseData[dto.RoleDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *RoleHandler) UpdateRolePermissionAction(ctx context.Context, i *dto.EditRolePermissionActions) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.roleUseCase.EditRolePermissionActions(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToUpdatePermissionAction"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditPermissionActionSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *RoleHandler) GetEntityActions(ctx context.Context, i *struct{ dto.AuthParams }) (
	*dto.EntityResponse[dto.ResultEntity[[]dto.EntityActionsDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.roleUseCase.GetEntityActions(req)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[[]dto.EntityActionsDto]]
	response.Body.Result = res
	return &response, nil
}

func (h *RoleHandler) GetRole(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.RoleDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.roleUseCase.GetRole(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.ROLE.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.RoleDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *RoleHandler) GetRoles(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.RoleDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.roleUseCase.GetRoles(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.ROLE.ID)
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.RoleDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *RoleHandler) GetRoleDefinitions(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.RoleActionDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.roleUseCase.GetRoleActions(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.RoleActionDto]]
	response.Body.PaginationResult = res
	return &response, nil
}
