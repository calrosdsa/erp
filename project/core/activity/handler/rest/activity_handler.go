package activity_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	activity_ucase "erp/project/core/activity/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ActivityHandler struct {
	sessionHelper   helpers.SessionHelper
	locale          helpers.Locale
	errorHelper     helpers.ErrorHelper
	activityUseCase activity_ucase.ActivityUseCase
	permission      repository.PermissionService
}

func NewActivityHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	activityUseCase activity_ucase.ActivityUseCase,
) {
	paths := NewActivityPaths(domain.ACTIVITY_BASE_ROUTE)
	tag := []string{"Activity"}
	handler := ActivityHandler{
		sessionHelper:   helpers.Session,
		locale:          helpers.Locale,
		errorHelper:     helpers.Error,
		activityUseCase: activityUseCase,
		permission:      permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "activity",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Activity",
		Tags:          tag,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateActivity)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-activity",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Activity",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditActivity)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-activity",
		Method:        http.MethodDelete,
		Path:          paths.Base,
		Summary:       "Delete  Activity",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.DeleteActivity)

	// huma.Register(api, huma.Operation{
	// 	OperationID:   "get-activities-party",
	// 	Method:        http.MethodGet,
	// 	Path:          paths.ByPartyID,
	// 	Summary:       "Get Activities Party",
	// 	Tags:          tag,
	// 	DefaultStatus: http.StatusOK,
	// 	Middlewares:   middlewares,
	// }, handler.GetActivitiesByPartyID)

}

// func (h *ActivityHandler) GetActivitiesByPartyID(ctx context.Context, i *dto.RequestEntity) (
//
//	*dto.ResponseData[[]dto.ActivityDto], error,
//
//	) {
//		req, err := h.sessionHelper.GetSession(ctx)
//		if err != nil {
//			return nil, huma.Error400BadRequest("Not Authorized", err)
//		}
//		// h.sessionHelper.AppendPaginationParams(req, i)
//		res, err := h.activityUseCase.GerActivitiesByPartyID(req, i)
//		if err != nil {
//			return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
//		}
//		var response dto.ResponseData[[]dto.ActivityDto]
//		response.Body.Result = res
//		return &response, err
//	}
func (h *ActivityHandler) DeleteActivity(ctx context.Context, i *dto.DeleteRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.activityUseCase.DeleteActivity(req, *i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToDelete")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.DeletedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ActivityHandler) EditActivity(ctx context.Context, i *dto.ActivityDataRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.activityUseCase.EditActivity(req, i.Body)
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

func (h *ActivityHandler) CreateActivity(ctx context.Context, i *dto.ActivityDataRequest) (
	*dto.ResponseData[dto.ActivityDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.activityUseCase.CreateActivity(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToAddComment")
	}
	var response dto.ResponseData[dto.ActivityDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CommentAddedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}
