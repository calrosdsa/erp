package profile_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	profile_ucase "erp/project/manage/profile/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ProfileHandler struct {
	sessionHelper helpers.SessionHelper
	usecase       profile_ucase.ProfileUseCase
	permission   repository.PermissionService
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
}

func NewProfileHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	usecase profile_ucase.ProfileUseCase,
	middlewares huma.Middlewares,
) {
	paths := NewProfilePaths(domain.PROFILE_ROUTE)
	tags := []string{"Profile"}
	handler := ProfileHandler{
		sessionHelper: helpers.Session,
		usecase:       usecase,
		permission:   permission,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "profiles",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Profiles",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetProfiles)

	huma.Register(api, huma.Operation{
		OperationID:   "get-profile-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get profile detail",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetProfileDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "get-profile-session",
		Method:        http.MethodGet,
		Path:          paths.Me,
		Summary:       "Get profile session",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetProfileSession)

	huma.Register(api, huma.Operation{
		OperationID:   "update-profile-session",
		Method:        http.MethodPut,
		Path:          paths.Me,
		Summary:       "Update profile session",
		Description:   "Update profile base on the currency session of the user",
		Tags:          tags,
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
	err = h.usecase.UpdateProfileSession(req, i)
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
	res, err := h.usecase.GetProfileSession(req)
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
	res, err := h.usecase.GetUserProfileDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(ctx, domain.USER.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.ProfileDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}



func (h *ProfileHandler) GetProfiles(ctx context.Context, i *dto.ProfilesRequest) (
	*dto.ResponseDataList[[]dto.ProfileDto], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.usecase.GetProfiles(req, *i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(ctx, domain.USER.ID)
	var response dto.ResponseDataList[[]dto.ProfileDto]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, err
}
