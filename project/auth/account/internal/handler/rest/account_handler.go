package rest_account

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	account_ucase "erp/project/auth/account/internal/usecase"
	"fmt"

	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type AccountHandler struct {
	jwtService     helpers.JwtHelper
	accountUcase   account_ucase.AccountUseCase
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	sessionService repository.SessionService
	errorHelper    helpers.ErrorHelper
}

func NewAccountHandler(
	api huma.API,
	helpers *helpers.Helpers,
	sessionService repository.SessionService,
	middlewares huma.Middlewares,
	accountUcase account_ucase.AccountUseCase,
) {
	paths := NewAccountPath(domain.ACCOUNT_BASE_ROUTE)
	tags := []string{"Account"}
	handler := AccountHandler{
		jwtService:     helpers.Jwt,
		accountUcase:   accountUcase,
		sessionHelper:  helpers.Session,
		sessionService: sessionService,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "sign-in",
		Method:        http.MethodPost,
		Path:          paths.SignIn,
		Summary:       "Sign in a user",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
	}, handler.SignIn)

	huma.Register(api, huma.Operation{
		OperationID:   "update-password",
		Method:        http.MethodPut,
		Path:          paths.Password,
		Summary:       "Update password",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdatePassword)

	huma.Register(api, huma.Operation{
		OperationID:   "get-account",
		Method:        http.MethodGet,
		Path:          paths.Account,
		Summary:       "Get user account",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAccount)

	huma.Register(api, huma.Operation{
		OperationID:   "get-sessions",
		Method:        http.MethodGet,
		Path:          paths.Sessions,
		Summary:       "Get user sessions",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetSessions)

	huma.Register(api, huma.Operation{
		OperationID:   "reset-password",
		Method:        http.MethodPost,
		Path:          paths.ResetPassoword,
		Summary:       "Reset Password",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
	}, handler.ResetPassword)
	huma.Register(api, huma.Operation{
		OperationID:   "change-password",
		Method:        http.MethodPost,
		Path:          paths.ChangePassword,
		Summary:       "Change Password",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
	}, handler.ChangePassword)
}
func (h *AccountHandler) GetSessions(ctx context.Context, i *struct{ dto.AuthParams }) (
	*dto.EntityResponse[[]dto.UserRelationDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.sessionService.GetUserRelations(req)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.EntityResponse[[]dto.UserRelationDto]
	response.Body.Result = res
	return &response, nil
}

func (h *AccountHandler) ChangePassword(ctx context.Context, i *dto.ChangePasswordRequest) (
	*dto.ResponseMessage, error) {
	lng := h.sessionHelper.ParseAcceptLanguage(i.AcceptLanguage)
	claims, err := h.jwtService.ExtractClaimsAdmin(i.Body.Token)
	fmt.Println("CLAIMS", claims)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(lng, err, "Error.FailToUpdate")
	}
	err = h.accountUcase.ChangePassword(ctx, i, claims)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(lng, err, "Error.FailToUpdate")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(lng),
	)
	return &response, nil
}

func (h *AccountHandler) ResetPassword(ctx context.Context, i *dto.ResetPasswordRequest) (*dto.ResponseMessage, error) {
	lng := h.sessionHelper.ParseAcceptLanguage(i.AcceptLanguage)
	err := h.accountUcase.ResetPassword(ctx, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToSendResetPasswordEmail"),
				helpers.OptionsLocale.WithLang(lng),
			), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.ResetPassword"),
		helpers.OptionsLocale.WithLang(lng),
		helpers.OptionsLocale.WithTemplate(map[string]string{
			"Email": i.Body.Email,
		}),
	)
	return &response, nil
}

func (h *AccountHandler) UpdatePassword(ctx context.Context, i *dto.UpdatePasswordRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	errorMsg, err := h.accountUcase.UpdatePassword(req, i)
	if err != nil {
		if errorMsg != nil {
			return nil, huma.Error400BadRequest(*errorMsg, err)
		}
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToChangePassword"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.ChangePasswordSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *AccountHandler) GetAccount(ctx context.Context, i *dto.AuthParams) (*dto.AccountResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	roleActions, err := h.sessionService.GetRoleActionsByRole(req, req.Role.ID)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToFetchData"),
				helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
			), err)
	}
	response := dto.AccountResponse{}
	response.Body.User = dto.UserDTOFromModel(&req.User)
	response.Body.Role = dto.RoleDTOFromModel(&req.Role)
	response.Body.Company = dto.CompanyDTOFromModel(&req.ActiveCompany)
	response.Body.Profile = dto.ProfileDTOFromModel(&req.Profile)
	response.Body.CompanyDefaults = dto.CompanyDefaultsDtoFromModel(&req.CompanyDefaults)
	response.Body.RoleActions = roleActions

	// plugins := h.pluginService.GetPlugins()
	// response.Body.AppConfig.Plugins = plugins
	return &response, nil
}

func (h *AccountHandler) SignIn(ctx context.Context, i *dto.SignInRequest) (*dto.SignInResponse, error) {
	// For now, we'll pass empty values since we can't access headers from context.Context in Huma handlers
	// In a full implementation, headers would be passed through middleware or request parameters
	ipAddress := ""
	userAgent := ""

	res, err := h.accountUcase.SignIn(ctx, i, ipAddress, userAgent)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid credentials", err)
	}	
	// session := common.GetSession(ctx)
	// fmt.Println(i.Body,session.LanguageCode)
	time.Sleep(time.Duration(1) * time.Second)
	return &res, nil
}
