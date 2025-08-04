package rest_company

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	company_ucase "erp/project/company/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type CompanyHandler struct {
	companyUseCase company_ucase.CompanyUseCase
	permission     repository.PermissionService
	sessionHelper  helpers.SessionHelper
	errorHelper    helpers.ErrorHelper
	locale         helpers.Locale
}

func NewCompanyHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	companyUseCase company_ucase.CompanyUseCase,
	middlewares huma.Middlewares,
) {
	paths := NewCompanyPath(domain.COMPANY_BASE_ROUTE)
	tags := []string{"Company"}
	handler := CompanyHandler{
		companyUseCase: companyUseCase,
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		permission:     permission,
		errorHelper:    helpers.Error,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "user-companies",
		Method:        http.MethodGet,
		Path:          paths.User,
		Summary:       "Get user companies",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAllUserCompanies)
	// huma.Register(api, huma.Operation{
	// 	OperationID:   "valid-parent-companies",
	// 	Method:        http.MethodGet,
	// 	Path:          paths.ValidParentCompanies,
	// 	Summary:       "Get Valid Parent Companies",
	// 	Tags:          tags,
	// 	DefaultStatus: http.StatusOK,
	// 	Middlewares:   middlewares,
	// }, handler.GetValidaParentCompanies)
	huma.Register(api, huma.Operation{
		OperationID:   "company-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get company Detial",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetCompanyDetail)
	// huma.Register(api, huma.Operation{
	// 	OperationID:   "get-company-by-uuid",
	// 	Method:        http.MethodGet,
	// 	Path:          paths.Uuid,
	// 	Summary:       "Get company by uuid",
	// 	Tags:          tags,
	// 	DefaultStatus: http.StatusOK,
	// }, handler.GetCompanyByUuid)

	huma.Register(api, huma.Operation{
		OperationID:   "create-company",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Company",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateCompany)

	huma.Register(api, huma.Operation{
		OperationID:   "company-account-setting",
		Method:        http.MethodGet,
		Path:          paths.AccountSetting,
		Summary:       "Company Account Setting",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.GetAccountSetting)
	huma.Register(api, huma.Operation{
		OperationID:   "edit-company-account-setting",
		Method:        http.MethodPut,
		Path:          paths.AccountSetting,
		Summary:       "Edit Company Account Setting",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.EditCompanyAccountSetting)
}

func (m *CompanyHandler) EditCompanyAccountSetting(ctx context.Context, d *dto.AccountSettingDataRequest,
	) (*dto.ResponseMessage,error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.companyUseCase.EditAccountSetting(req, d.Body)
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

func (h *CompanyHandler) GetAccountSetting(ctx context.Context, i *dto.RequestData) (
	*dto.EntityResponse[dto.AccountSettingsDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.companyUseCase.GetAccountSetting(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.AccountSettingsDto]
	response.Body.Result = res
	return &response, err
}

func (h *CompanyHandler) CreateCompany(ctx context.Context, i *dto.CreateCompanyRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.companyUseCase.CreateCompany(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateCompany")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateCompanySuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

// func (h *CompanyHandler) GetCompanyByUuid(ctx context.Context, i *dto.UuidRequest) (
// 	*dto.EntityResponse[model.Company], error) {
// 	res, err := h.companyUseCase.GetCompanyByUuid(ctx, i.Uuid)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest(h.locale.MustLocalize(
// 			helpers.OptionsLocale.WithID("Error.FailToFetchData"),
// 			helpers.OptionsLocale.WithLang(i.AcceptLanguage),
// 		), err)
// 	}
// 	var response dto.EntityResponse[model.Company]
// 	response.Body.Result = res
// 	return &response, err
// }

func (h *CompanyHandler) GetCompanyDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.CompanyDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.companyUseCase.GetCompanyDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.CompanyDto]]
	response.Body.Result = res
	response.Body.AssociatedActions = h.getActions(ctx)
	return &response, err
}
func (h *CompanyHandler) GetAllUserCompanies(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.companyUseCase.GetAllUserCompanies(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.COMPANY.ID)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

// func (h *CompanyHandler) GetCompanies(ctx context.Context, i *dto.RequestPaginationData) (
// 	*dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]], error,
// ) {
// 	req, err := h.sessionHelper.GetSession(ctx)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Not Authorized", err)
// 	}
// 	h.sessionHelper.AppendPaginationParams(req, i)
// 	res, err := h.companyUseCase.GetAllUserCompanies(req, i)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, huma.Error400BadRequest("Error", err)
// 	}
// 	actions := h.permission.GetActions(req.Ctx, domain.COMPANY.ID)
// 	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]]{}
// 	response.Body.PaginationResult = res
// 	response.Body.Actions = actions
// 	return response, err
// }

func (h *CompanyHandler) GetValidaParentCompanies(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.companyUseCase.GetAllUserCompanies(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response *dto.PaginationResponse[dto.PaginationResult[[]dto.CompanyDto]]
	response.Body.PaginationResult = res
	return response, err
}

func (h *CompanyHandler) getActions(ctx context.Context) map[int][]dto.ActionDto {
	var ids []int64
	ids = append(ids, domain.COMPANY.ID,domain.LEDGER.ID)
	r := h.permission.GetEntitiesActions(ctx, ids)
	return r
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
