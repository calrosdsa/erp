package company

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/domain"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/internal/app/service/services/account_service"
	"erp/internal/app/service/services/company_service"
	"fmt"

	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type CompanyHandler struct {
	companyService *company_service.CompanyService
	roleService *account_service.RoleService
	sessionHelper  helpers.SessionHelper
	errorHelper helpers.ErrorHelper
	locale         helpers.Locale
}

func NewHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag string,
	middlewares huma.Middlewares,
) {
	paths := NewCompanyPath(base)
	handler := CompanyHandler{
		companyService: services.CompanyService,
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		roleService: services.RoleService,
		errorHelper: helpers.Error,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "companies",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get User companies",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetCompanies)
	huma.Register(*api, huma.Operation{
		OperationID:   "user-companies",
		Method:        http.MethodGet,
		Path:          paths.User,
		Summary:       "Get user companies",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAllUserCompanies)
	huma.Register(*api, huma.Operation{
		OperationID:   "valid-parent-companies",
		Method:        http.MethodGet,
		Path:          paths.ValidParentCompanies,
		Summary:       "Get Valid Parent Companies",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetValidaParentCompanies)
	huma.Register(*api, huma.Operation{
		OperationID:   "company-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get company Detial",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetCompanyDetail)
	huma.Register(*api, huma.Operation{
		OperationID:   "get-company-by-uuid",
		Method:        http.MethodGet,
		Path:          paths.Uuid,
		Summary:       "Get company by uuid",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
	}, handler.GetCompanyByUuid)

	huma.Register(*api, huma.Operation{
		OperationID:   "create-company",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Company",
		Tags:          []string{tag},
		DefaultStatus: http.StatusCreated,
		Middlewares: middlewares,
	}, handler.CreateCompany)
}

func (h *CompanyHandler) CreateCompany(ctx context.Context, i *dto.CreateCompanyRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.companyService.CreateCompany(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToCreateCompany"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateCompanySuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *CompanyHandler) GetCompanyByUuid(ctx context.Context, i *dto.UuidRequest) (
	*dto.EntityResponse[model.Company], error) {
	res, err := h.companyService.GetCompanyByUuid(ctx, i.Uuid)
	if err != nil {
		return nil, huma.Error400BadRequest(h.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Error.FailToFetchData"),
			helpers.OptionsLocale.WithLang(i.AcceptLanguage),
		), err)
	}
	var response dto.EntityResponse[model.Company]
	response.Body.Result = res
	return &response, err
}

func (h *CompanyHandler) GetCompanyDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.CompanyDto], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.companyService.GetCompanyDetail(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	return &res, err
}

func (h *CompanyHandler) GetAllUserCompanies(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.companyService.GetAllUserCompanies(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
	}
	actions := h.roleService.GetActions(req,domain.COMPANY)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *CompanyHandler) GetCompanies(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.companyService.GetCompanies(req, i)
	if err != nil {
		fmt.Println(err)
		return nil, huma.Error400BadRequest("Error", err)
	}
	actions := h.roleService.GetActions(req,domain.COMPANY)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *CompanyHandler) GetValidaParentCompanies(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.companyService.GetValidaParentCompanies(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest("Error", err)
	}
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]]{}
	response.Body.PaginationResult = res
	return response, err
}

// func (h CompanyHandler) SignIn(ctx echo.Context) error {
// 	token ,err := h.jwtService.GenerateToken(common.Claims{
// 		ID: 1,
// 		Uuid: "5f4a8a86-0fd0-4c71-8615-f44a8475d595",
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	time.Sleep(time.Duration(2)* time.Second)
// 	return ctx.JSON(http.StatusOK, SignInResponse{
// 		AccessToken: token,
// 	})
// }
