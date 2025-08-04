package account_api

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/internal/app/service/services/account_service"
	"erp/internal/app/service/services/jwt_service"
	pluginservice "erp/internal/app/service/services/plugin_service"

	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type AccountHandler struct {
	jwtService     *jwt_service.JwtService
	accountService *account_service.AccountService
	sessionHelper  helpers.SessionHelper
	pluginService  *pluginservice.PluginService
	locale         helpers.Locale
	sessionService *account_service.SessionService
	roleService *account_service.RoleService
}

func NewHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag string,
	middlewares huma.Middlewares,
) {
	paths := NewAccountPath(base)
	handler := AccountHandler{
		jwtService:     services.JwtService,
		accountService: services.AccountService,
		sessionHelper:  helpers.Session,
		sessionService: services.SessionService,
		pluginService:  services.PluginService,
		locale:         helpers.Locale,
		roleService: services.RoleService,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "sign-in",
		Method:        http.MethodPost,
		Path:          paths.SignIn,
		Summary:       "Sign in a user",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
	}, handler.SignIn)

	huma.Register(*api, huma.Operation{
		OperationID:   "update-password",
		Method:        http.MethodPut,
		Path:          paths.Password,
		Summary:       "Update password",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdatePassword)

	huma.Register(*api, huma.Operation{
		OperationID:   "get-account",
		Method:        http.MethodGet,
		Path:          paths.Account,
		Summary:       "Get user account",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAccount)

	huma.Register(*api, huma.Operation{
		OperationID:   "get-sessions",
		Method:        http.MethodGet,
		Path:          paths.Sessions,
		Summary:       "Get user sessions",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetSessions)

}
func (h *AccountHandler) GetSessions(ctx context.Context, i *struct{ dto.AuthParams }) (
	*dto.EntityResponse[[]dto.UserRelationDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
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

func (h *AccountHandler) ResetPassword(ctx context.Context, i *dto.ResetPasswordRequest) (*dto.ResponseMessage, error) {
	err := h.accountService.ResetPassword(ctx, i)
	if err != nil {
		return nil, huma.Error400BadRequest(
			h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Error.FailToSendResetPasswordEmail"),
				helpers.OptionsLocale.WithLang(""),
			), err)
	}
	var response dto.ResponseMessage
	// response.Body.Message = h.locale.MustLocalize(
	// 	helpers.OptionsLocale.WithID(""),
	// )
	return &response, nil
}

func (h *AccountHandler) UpdatePassword(ctx context.Context, i *dto.UpdatePasswordRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	errorMsg, err := h.accountService.UpdatePassword(req, i)
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
	roleActions,err := h.roleService.GetRoleActionsByRole(req,req.Role.ID)
	response := dto.AccountResponse{}
	response.Body.User = dto.UserDTOFromModel(&req.User)
	response.Body.Role = dto.RoleDTOFromModel(&req.Role)
	response.Body.Company = dto.CompanyDTOFromModel(&req.ActiveCompany)
	response.Body.Profile = dto.ProfileDTOFromModel(&req.Profile)	
	response.Body.RoleActions = roleActions

	// plugins := h.pluginService.GetPlugins()
	// response.Body.AppConfig.Plugins = plugins
	return &response, nil
}

func (h *AccountHandler) SignIn(ctx context.Context, i *dto.SignInRequest) (*dto.SignInResponse, error) {
	res, err := h.accountService.SignIn(ctx, i)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid credentials", err)
	}
	token, err := h.jwtService.GenerateToken(common.Claims{
		ID:   res.Body.U.ID,
		Uuid: res.Body.U.UUID,
	})
	if err != nil {
		return nil, huma.Error400BadRequest("Fial to generate token", err)
	}
	// session := common.GetSession(ctx)
	// fmt.Println(i.Body,session.LanguageCode)
	time.Sleep(time.Duration(1) * time.Second)
	res.Body.AccessToken = token
	return &res, nil
}

// func (h AccountHandler) SignIn(ctx echo.Context) error {
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
