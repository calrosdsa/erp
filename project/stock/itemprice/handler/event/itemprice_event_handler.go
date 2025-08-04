package itemprice_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	itemprice_repo "erp/project/stock/itemprice/repository"
)

type itemPriceEventHandler struct {
	eventRepo itemprice_repo.ItemPriceEventRepo
	emitLog logger.EmitLog
}

func NewEventHandler(
	logger logger.Logger,
	eventRepo itemprice_repo.ItemPriceEventRepo,
	bus bus.Bus,
){
	h := itemPriceEventHandler{
		eventRepo: eventRepo,
		emitLog: logger.EmitLog("item-price-event-handler"),
	}
	bus.RegisterHandler(domain.ItemCreatedEvent,h.OnItemCreated())
	bus.RegisterHandler(domain.ItemEditedEvent,h.OnItemEdited())
}
func (h *itemPriceEventHandler)OnItemEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnItemEdited"))
				}
			}()
			payload,ok := e.Data.(event.ItemCreatedEventData)
			//Assign current context  to the Req object
			payload.Req.Ctx = ctx
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.eventRepo.OnItemEdited(ctx,payload)
			return 
		},
		AbortOnError: true,
		Matcher:domain.ItemEditedEvent,
	}
}

func (h *itemPriceEventHandler)OnItemCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnItemCreated"))
				}
			}()
			payload,ok := e.Data.(event.ItemCreatedEventData)
			//Assign current context  to the Req object
			payload.Req.Ctx = ctx
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.eventRepo.OnItemCreated(ctx,payload)
			return 
		},
		AbortOnError: true,
		Matcher:domain.ItemCreatedEvent,
	}
}