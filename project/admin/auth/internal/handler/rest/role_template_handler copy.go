package auth_admin_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	auth_admin_ucase "erp/project/admin/auth/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type RoleTemplateHandler struct {
	jwt               helpers.JwtHelper
	sessionHelper     helpers.SessionHelper
	locale            helpers.Locale
	errorHelper       helpers.ErrorHelper
	roleTemplateUcase auth_admin_ucase.RoleTemplateUseCase
}

func NewRoleTemplateHandler(
	api huma.API,
	helpers *helpers.Helpers,
	roleTemplateUcase auth_admin_ucase.RoleTemplateUseCase,
	middlewares huma.Middlewares,
) {
	base := domain.ROLE_TEMPLATE_BASE_ROUTE
	tags := []string{"Role Template"}
	path := NewRoleTemplatePaths(base)
	h := RoleTemplateHandler{
		jwt:               helpers.Jwt,
		sessionHelper:     helpers.Session,
		locale:            helpers.Locale,
		errorHelper:       helpers.Error,
		roleTemplateUcase: roleTemplateUcase,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "role-templates",
		Method:        http.MethodGet,
		Summary:       "Role Templates",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetRoleTemplates)

	huma.Register(api, huma.Operation{
		OperationID:   "create-role-template",
		Method:        http.MethodPost,
		Summary:       "Create Role Template",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateRoleTemplate)

	huma.Register(api, huma.Operation{
		OperationID:   "role-template",
		Method:        http.MethodGet,
		Summary:       "Role Template",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetRoleTemplate)
}

func (h *RoleTemplateHandler) CreateRoleTemplate(ctx context.Context, d *dto.CreateRoleTemplateRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.roleTemplateUcase.CreateRoleTemplate(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err,"Error.FailToCreate")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *RoleTemplateHandler) GetRoleTemplate(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.RoleTemplateDto]], error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.roleTemplateUcase.GetRoleTemplate(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.RoleTemplateDto]]
	response.Body.Result = res
	return &response, nil
}

func (h *RoleTemplateHandler) GetRoleTemplates(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.RoleTemplateDto]], error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.roleTemplateUcase.GetRoleTemplates(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.RoleTemplateDto]]
	response.Body.PaginationResult = res
	return &response, nil
}
