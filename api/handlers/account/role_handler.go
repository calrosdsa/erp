package account_api

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/internal/app/service/services/account_service"
	"erp/internal/app/service/services/jwt_service"
	"erp/internal/domain"

	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type RoleHandler struct {
	jwtService    *jwt_service.JwtService
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	roleService   *account_service.RoleService
	errorHelper   helpers.ErrorHelper
}

func NewRoleHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag string,
	middlewares huma.Middlewares,
) {
	paths := NewRolePaths(base)
	handler := RoleHandler{
		jwtService:    services.JwtService,
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		roleService:   services.RoleService,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "get roles",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get roles",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetRoles)
	huma.Register(*api, huma.Operation{
		OperationID:   "get role definitions",
		Method:        http.MethodGet,
		Path:          paths.RoleDefinitions,
		Summary:       "Get role definitions",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetRoleDefinitions)

	huma.Register(*api, huma.Operation{
		OperationID:   "get entity actions",
		Method:        http.MethodGet,
		Path:          paths.EntityActions,
		Summary:       "Get entity actions",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetEntityActions)

	huma.Register(*api, huma.Operation{
		OperationID:   "get role",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get role",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetRole)

	huma.Register(*api, huma.Operation{
		OperationID:   "update role permision action",
		Method:        http.MethodPost,
		Path:          paths.PermissionActions,
		Summary:       "Update Permission Actions",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateRolePermissionAction)

	huma.Register(*api, huma.Operation{
		OperationID:   "create role",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Role",
		Tags:          []string{tag},
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateRole)
}

func (h *RoleHandler) CreateRole(ctx context.Context, i *dto.RoleRequestData) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.roleService.CreateRole(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateRole")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateRoleSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *RoleHandler) UpdateRolePermissionAction(ctx context.Context, i *dto.EditRolePermissionActions) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.roleService.EditRolePermissionActions(req, i)
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
	res, err := h.roleService.GetEntityActions(req)
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
	res, err := h.roleService.GetRole(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	actions := h.roleService.GetActions(req, domain.ROLE)
	var response dto.EntityResponse[dto.ResultEntity[dto.RoleDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *RoleHandler) GetRoles(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.RoleDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.roleService.GetRoles(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	actions := h.roleService.GetActions(req, domain.ROLE)
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.RoleDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *RoleHandler) GetRoleDefinitions(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.RoleActionDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.roleService.GetRoleActions(req, i)
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
