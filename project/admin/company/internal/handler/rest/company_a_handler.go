package a_company_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	a_company_ucase "erp/project/admin/company/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type CompanyAHandler struct {
	jwt           helpers.JwtHelper
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	companyAUCase a_company_ucase.AdminCompanyUCase
}

func NewCompanyAHandler(
	api huma.API,
	helpers *helpers.Helpers,
	companyAUCase a_company_ucase.AdminCompanyUCase,
	middlewares huma.Middlewares,
) {
	base := domain.A_COMPANY_BASE_ROUTE
	tags := []string{"Admin Company"}
	path := NewACompanyAdminPaths(base)
	h := CompanyAHandler{
		jwt:           helpers.Jwt,
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		companyAUCase: companyAUCase,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "a-companies",
		Method:        http.MethodGet,
		Summary:       "Admin Companies",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCompanies)

	huma.Register(api, huma.Operation{
		OperationID:   "a-create-company",
		Method:        http.MethodPost,
		Summary:       "Create Company Adm",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateCompany)

	huma.Register(api, huma.Operation{
		OperationID:   "a-company-detail",
		Method:        http.MethodGet,
		Summary:       "Admin Company Detail",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCompany)

	huma.Register(api, huma.Operation{
		OperationID:   "a-company-modules",
		Method:        http.MethodGet,
		Summary:       "Admin Company Modules",
		Tags:          tags,
		Path:          path.Modules,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCompanyModules)

	huma.Register(api, huma.Operation{
		OperationID:   "a-add-company-modules",
		Method:        http.MethodPost,
		Summary:       "Add Company Modules",
		Tags:          tags,
		Path:          path.Modules,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.AddCompanyModules)

	huma.Register(api, huma.Operation{
		OperationID:   "a-company-users",
		Method:        http.MethodGet,
		Summary:       "Get Company Users",
		Tags:          tags,
		Path:          path.User,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCompanyUsers)

	huma.Register(api, huma.Operation{
		OperationID:   "add-company-user",
		Method:        http.MethodPost,
		Summary:       "Add Company User",
		Tags:          tags,
		Path:          path.User,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.AddCompanyUser)
}
func (h *CompanyAHandler) AddCompanyUser(ctx context.Context, d *dto.CreateUserAdminRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.companyAUCase.AddCompanyUser(req, d)
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

func (h *CompanyAHandler) GetCompanyUsers(ctx context.Context, d *dto.RequestData) (
	*dto.EntityResponse[[]dto.UserDto], error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.companyAUCase.GetCompanyUsers(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.UserDto]
	response.Body.Result = res
	return &response, nil
}

func (h *CompanyAHandler) AddCompanyModules(ctx context.Context, d *dto.AddCompanyModules) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.companyAUCase.AddCompanyModules(req, d)
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

func (h *CompanyAHandler) GetCompanyModules(ctx context.Context, d *dto.RequestData) (
	*dto.EntityResponse[[]dto.CompanyEntityDto], error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.companyAUCase.GetCompanyModules(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[[]dto.CompanyEntityDto]
	response.Body.Result = res
	return &response, nil
}

func (h *CompanyAHandler) GetCompany(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.CompanyDto]], error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.companyAUCase.GetCompany(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.CompanyDto]]
	response.Body.Result = res
	return &response, nil
}

func (h *CompanyAHandler) GetCompanies(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]], error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.companyAUCase.GetParentCompanies(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]]
	response.Body.PaginationResult = res
	return &response, nil
}


func (h *CompanyAHandler) CreateCompany(ctx context.Context, d *dto.CreateCompanyAdminRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.companyAUCase.CreateCompany(req, d)
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
