package eventbooking_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	event_ucase "erp/project/regate/event/internal/usecase"
	regate_domain "erp/project/regate/internal/domain"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type EventBookingHandler struct {
	sessionHelper     helpers.SessionHelper
	locale            helpers.Locale
	errorHelper       helpers.ErrorHelper
	eventBookingUcase event_ucase.EventBookingUseCase
	permission        repository.PermissionService
}

func NewEventBookingHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	eventBookingUcase event_ucase.EventBookingUseCase,
) {
	paths := NewEventBookingPaths(domain.EVENT_BASE_ROUTE)
	tag := []string{"Event"}
	handler := EventBookingHandler{
		sessionHelper:     helpers.Session,
		locale:            helpers.Locale,
		errorHelper:       helpers.Error,
		eventBookingUcase: eventBookingUcase,
		permission:        permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "create-event-booking",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Event",
		Tags:          tag,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateEventBooking)

	huma.Register(api, huma.Operation{
		OperationID:   "get-event-bookings",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Event Booking",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetEventBookings)
	huma.Register(api, huma.Operation{
		OperationID:   "get-event-booking",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Event Booking",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetEventBooking)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-event",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit event",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditEvent)

	huma.Register(api, huma.Operation{
		OperationID:   "event-update-status",
		Method:        http.MethodPut,
		Path:          paths.UpdateStatus,
		Summary:       "Event Update Status",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-event-batch",
		Method:        http.MethodDelete,
		Path:          paths.DeleteInBatch,
		Summary:       "Delete Event Batch",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.DeleteEventBatch)
}
func (h *EventBookingHandler) DeleteEventBatch(ctx context.Context, i *dto.DeleteEventBatchRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.eventBookingUcase.DeleteEventBatch(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToDelete")
		// return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.DeletedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *EventBookingHandler) UpdateStatus(ctx context.Context, i *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.eventBookingUcase.UpdateStatus(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdate")
		// return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *EventBookingHandler) EditEvent(ctx context.Context, d *dto.EventBookingDataRequest) (*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.eventBookingUcase.EditEvent(req, d.Body)
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

func (h *EventBookingHandler) GetEventBooking(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.EventBookingDetail]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.eventBookingUcase.GetEventBooking(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, regate_domain.EVENT_BOOKING.ID)

	response := &dto.EntityResponse[dto.ResultEntity[dto.EventBookingDetail]]{}
	response.Body.Result = res
	response.Body.Actions = actions
	return response, err
}

func (h *EventBookingHandler) GetEventBookings(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.EventBookingDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.eventBookingUcase.GetEventBookings(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, regate_domain.EVENT_BOOKING.ID)

	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.EventBookingDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *EventBookingHandler) CreateEventBooking(ctx context.Context, i *dto.EventBookingDataRequest) (
	*dto.ResponseData[dto.EventBookingDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.eventBookingUcase.CreateEventBooking(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.EventBookingDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}
