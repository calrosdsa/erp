package item_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	item_repo "erp/project/stock/item/repository"
)

type itemPriceEventHandler struct {
	repo    item_repo.ItemEventRepository
	emitLog logger.EmitLog
}

func NewEventHandler(
	logger logger.Logger,
	repo item_repo.ItemEventRepository,
	bus bus.Bus,
) {
	h := itemPriceEventHandler{
		// eventRepo: eventRepo,

		emitLog:   logger.EmitLog("item-price-event-handler"),
	}
	bus.RegisterHandler(domain.OrderCreatedEvent, h.OnPurchaseOrderCreated())
}

func (h *itemPriceEventHandler) OnPurchaseOrderCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnPurchaseOrderCreated"))
				}
			}()
			payload, ok := e.Data.(event.OrderEventData)
			//Assign current context  to the Req object
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.repo.OnPurchaseOrderCreated(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.OrderCreatedEvent,
	}
}