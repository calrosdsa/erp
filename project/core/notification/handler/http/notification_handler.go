package notification_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	notification_ucase "erp/project/core/notification/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type NotifiationHandler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	usecase       notification_ucase.NotificationUcase
}

func NewNotificationHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	usecase notification_ucase.NotificationUcase,
) {
	base := domain.NOTIFICATION_BASE_ROUTE
	tags := []string{"Notification"}
	h := NotifiationHandler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		usecase:       usecase,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "notification-count",
		Method:        http.MethodGet,
		Summary:       "Notification Count",
		Tags:          tags,
		Path:          base + "/count",
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetNotificationCount)

	huma.Register(api, huma.Operation{
		OperationID:   "notification",
		Method:        http.MethodGet,
		Summary:       "Notification",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetNotifications)
}

func (h *NotifiationHandler) GetNotificationCount(ctx context.Context,d *struct{}) (
	*dto.NotificationCountDto,error,
){
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.usecase.NotifiationCount(req)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.NotificationCountDto
	response.Body.Count = res
	return &response, nil
}

func (h *NotifiationHandler) GetNotifications(ctx context.Context,d *dto.NotificationsRequest) (
	*dto.ResponseDataList[[]dto.NotificationDto],error,
){
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.usecase.GetNotifications(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	return &res, nil
}