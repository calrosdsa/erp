package item_inventory_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	item_inventory_repo "erp/project/stock/item_inventory/internal/repository"
)

type itemInventoryEventHandler struct {
	emitLog logger.EmitLog
	itemInventoryRepo item_inventory_repo.ItemInventoryRepo
}

func NewItemInventoryEvent(
	logger logger.Logger,
	itemInventoryRepo item_inventory_repo.ItemInventoryRepo,
	bus bus.Bus,
) {
	h := itemInventoryEventHandler{
		itemInventoryRepo:itemInventoryRepo,
		emitLog:logger.EmitLog("item-inventory-event"),
	}
	bus.RegisterHandler(domain.ItemCreatedEvent,h.OnItemCreated())
	bus.RegisterHandler(domain.ItemEditedEvent,h.OnItemEdited())
}

func (h *itemInventoryEventHandler) OnItemEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnItemEdited"))
				}
			}()
			payload,ok := e.Data.(event.ItemCreatedEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemInventoryRepo.OnEditInventory(ctx,payload)
			return 
		},
		AbortOnError: true,
		Matcher: domain.ItemEditedEvent,
	}
}

func (h *itemInventoryEventHandler) OnItemCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnItemCreated"))
				}
			}()
			payload,ok := e.Data.(event.ItemCreatedEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemInventoryRepo.OnCreateItemInventory(ctx,payload)
			return 
		},
		AbortOnError: true,
		Matcher: domain.ItemCreatedEvent,
	}
}