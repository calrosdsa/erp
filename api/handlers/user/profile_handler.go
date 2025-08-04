package user

import (
	"context"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/internal/app/service/services/account_service"
	userservice "erp/internal/app/service/services/user_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ProfileHandler struct {
	sessionHelper  helpers.SessionHelper
	profileService *userservice.ProfileService
	roleService    *account_service.RoleService
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
}

func NewProfileHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag []string,
	middlewares huma.Middlewares,
) {
	paths := NewProfilePaths(base)
	handler := ProfileHandler{
		sessionHelper:  helpers.Session,
		profileService: services.ProfileService,
		roleService:    services.RoleService,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "get company user profiles",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get company profile users",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetCompanyUserProfiles)

	huma.Register(*api, huma.Operation{
		OperationID:   "get profile detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get profile detail",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetProfileDetail)
	huma.Register(*api, huma.Operation{
		OperationID:   "create user",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create User",
		Tags:          tag,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateUser)

	huma.Register(*api, huma.Operation{
		OperationID:   "get-profile-session",
		Method:        http.MethodGet,
		Path:          paths.Me,
		Summary:       "Get profile session",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetProfileSession)

	huma.Register(*api, huma.Operation{
		OperationID:   "update-profile-session",
		Method:        http.MethodPut,
		Path:          paths.Me,
		Summary:       "Update profile session",
		Description:   "Update profile base on the currency session of the user",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateProfileSession)
}
func (h *ProfileHandler) UpdateProfileSession(ctx context.Context, i *dto.UpdateProfileRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.profileService.UpdateProfileSession(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdateProfile")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdateProfileSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)

	return &response, nil
}

func (h *ProfileHandler) GetProfileSession(ctx context.Context, i *struct {
	dto.AuthParams
}) (
	*dto.EntityResponse[dto.ResultEntity[dto.ProfileDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.profileService.GetProfileSession(req)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.ProfileDto]]
	response.Body.Result = res
	return &response, nil
}

func (h *ProfileHandler) GetProfileDetail(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ProfileDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.profileService.GetUserProfileDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.roleService.GetActions(req, domain.USER)
	var response dto.EntityResponse[dto.ResultEntity[dto.ProfileDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *ProfileHandler) CreateUser(ctx context.Context, i *dto.CreateUserRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.profileService.CreateUser(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateUser")
	}
	var res dto.ResponseMessage
	res.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateUserSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &res, nil
}

func (h *ProfileHandler) GetCompanyUserProfiles(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.ProfileL]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.profileService.GetCompanyUserProfiles(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ActionNoAllowed")
	}
	actions := h.roleService.GetActions(req, domain.USER)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.ProfileL]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}
