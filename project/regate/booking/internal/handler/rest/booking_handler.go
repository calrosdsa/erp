package booking_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	booking_ucase "erp/project/regate/booking/internal/usecase"
	regate_domain "erp/project/regate/internal/domain"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type BookingHandler struct {
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	bookingUseCase booking_ucase.BookingUseCase
	permission     repository.PermissionService
}

func NewBookingHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	bookingUseCase booking_ucase.BookingUseCase,
) {
	paths := NewBookingPaths(domain.BOOKING_BASE_ROUTE)
	tag := []string{"Booking"}
	handler := BookingHandler{
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
		bookingUseCase: bookingUseCase,
		permission:     permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "create booking",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Booking",
		Tags:          tag,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateBooking)

	huma.Register(api, huma.Operation{
		OperationID:   "get-bookings",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Bookings",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetBookings)
	huma.Register(api, huma.Operation{
		OperationID:   "get-booking",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Booking",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetBooking)

	huma.Register(api, huma.Operation{
		OperationID:   "validate-booking",
		Method:        http.MethodPost,
		Path:          paths.Validate,
		Summary:       "Validate Booking",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.ValidateBooking)

	huma.Register(api, huma.Operation{
		OperationID:   "update-booking-status",
		Method:        http.MethodPut,
		Path:          paths.UpdateStatus,
		Summary:       "Update Booking Status",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateBookingStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "update-paid-amount",
		Method:        http.MethodPut,
		Path:          paths.UpdatePaidAmount,
		Summary:       "Update Paid Amount",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditAmountPaid)

	huma.Register(api, huma.Operation{
		OperationID:   "booking-reschedule",
		Method:        http.MethodPut,
		Path:          paths.Reschedule,
		Summary:       "Booking reschedule",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.BookingReschedule)

	huma.Register(api, huma.Operation{
		OperationID:   "update-booking-batch",
		Method:        http.MethodPost,
		Path:          paths.UpdateBookingBatch,
		Summary:       "Update Booking Batch",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateBookingBatch)
}
func (h *BookingHandler) UpdateBookingBatch(ctx context.Context, i *dto.UpdateBookingBatchRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.bookingUseCase.UpdateBookingBatch(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *BookingHandler) BookingReschedule(ctx context.Context, i *dto.BookingRescheduleRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.bookingUseCase.BookingReschedule(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *BookingHandler) EditAmountPaid(ctx context.Context, i *dto.BookingPaymentRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.bookingUseCase.EditBookingPaidAmount(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *BookingHandler) UpdateBookingStatus(ctx context.Context, i *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.bookingUseCase.UpdateBookingState(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *BookingHandler) ValidateBooking(ctx context.Context, i *dto.ValidateBookingRequest) (
	*dto.ResponseData[dto.ValidateBookingData], error,
) {
	fmt.Println("VALIDATE BOOKING...")
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.bookingUseCase.ValidateBooking(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[dto.ValidateBookingData]
	response.Body.Result = res
	return &response, err
}

func (h *BookingHandler) GetBooking(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.BookingDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.bookingUseCase.GetBooking(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, regate_domain.BOOKING.ID)

	response := &dto.EntityResponse[dto.ResultEntity[dto.BookingDto]]{}
	response.Body.Result = res
	response.Body.Actions = actions
	return response, err
}

func (h *BookingHandler) GetBookings(ctx context.Context, i *dto.RequestBookings) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.BookingDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.bookingUseCase.GetBookings(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, regate_domain.BOOKING.ID)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.BookingDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *BookingHandler) CreateBooking(ctx context.Context, i *dto.CreateBookingRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.bookingUseCase.CreateBooking(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}
