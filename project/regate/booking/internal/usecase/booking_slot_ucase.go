package booking_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	booking_repo "erp/project/regate/booking/internal/repository"
)

type BookingSlotUseCase interface {
	GetBookingSlots(req *common.RequestContext, d *dto.RequestBookingSlots) (
		res dto.BookingScheduleBody, err error)
}

type bookingSlotUseCase struct {
	bookingSlotRepo booking_repo.BookingSlotRepository
	emitLog logger.EmitLog
}

func NewBookingSlotUseCase(
	logger logger.Logger,
	bookingSlotRepo booking_repo.BookingSlotRepository,
)BookingSlotUseCase{
	return &bookingSlotUseCase{
		emitLog: logger.EmitLog("booking-slot-usecase"),
		bookingSlotRepo: bookingSlotRepo,
	}
}

func (u *bookingSlotUseCase)GetBookingSlots(req *common.RequestContext, d *dto.RequestBookingSlots) (
	res dto.BookingScheduleBody, err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetBookingSlots"))
		}
	}()
	res,err = u.bookingSlotRepo.GetBookingSlots(req,d)
	return
}