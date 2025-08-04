package booking_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	booking_ucase "erp/project/regate/booking/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type BookingSlotHandler struct {
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	bookingSlotUseCase booking_ucase.BookingSlotUseCase
	permission     repository.PermissionService
}

func NewBookingSlotHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	bookingSlotUseCase booking_ucase.BookingSlotUseCase,
) {
	paths := NewBookingSlotPaths(domain.BOOKING_SLOT_BASE_ROUTE)
	tag := []string{"Booking Slot"}
	handler := BookingSlotHandler{
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
		bookingSlotUseCase: bookingSlotUseCase,
		permission:     permission,
	}

	
	huma.Register(api, huma.Operation{
		OperationID:   "booking-slots",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Booking Slots",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetBookingSlots)
	
}


func (h *BookingSlotHandler) GetBookingSlots(ctx context.Context, i *dto.RequestBookingSlots) (
	*dto.BookingScheduleBody, error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.bookingSlotUseCase.GetBookingSlots(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	return &res, err
}
