package order_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	order_repo "erp/project/order/repository"
	"fmt"
)

type orderEventHandler struct {
	bus            bus.Bus
	orderEventRepo order_repo.OrderEventRepository
	emitLog        logger.EmitLog
}

func NewOrderEventHandler(
	bus bus.Bus,
	logger logger.Logger,
	orderEventRepo order_repo.OrderEventRepository,
) {
	handler := orderEventHandler{
		bus:            bus,
		emitLog:        logger.EmitLog("order-event"),
		orderEventRepo: orderEventRepo,
	}
	bus.RegisterHandler(domain.ReceiptCreatedEvent, handler.OnReceiptCreated())
	bus.RegisterHandler(domain.ReceiptCancelledEvent, handler.OnReceiptCancelled())
	bus.RegisterHandler(domain.ReceiptSubmittedEvent, handler.OnReceiptSubmitted())
	bus.RegisterHandler(domain.InvoiceSubmittedEvent, handler.OnInvoiceSubmitted())
	bus.RegisterHandler(domain.InvoiceCancelledEvent, handler.OnInvoiceCancelled())
}

func (h *orderEventHandler) OnInvoiceCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.orderEventRepo.OnInvoiceCancelled(ctx, payload)
			return err
		},
		Matcher:      domain.InvoiceCancelledEvent,
		AbortOnError: false,
	}
}

func (h *orderEventHandler) OnInvoiceSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.StatusInvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.orderEventRepo.OnInvoiceSubmitted(ctx, payload)
			return err
		},
		Matcher:      domain.InvoiceSubmittedEvent,
		AbortOnError: false,
	}
}

func (h *orderEventHandler) OnReceiptCancelled() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.StatusReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.orderEventRepo.OnReceiptCancelled(ctx, payload)
			return err
		},
		Matcher:      domain.ReceiptCancelledEvent,
		AbortOnError: true,
	}
}

func (h *orderEventHandler) OnReceiptSubmitted() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			payload, ok := e.Data.(event.StatusReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err := h.orderEventRepo.OnReceiptSubmitted(ctx, payload)
			return err
		},
		Matcher:      domain.ReceiptSubmittedEvent,
		AbortOnError: true,
	}
}

func (h *orderEventHandler) OnReceiptCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			fmt.Println("ORDER RECEIPT EVENT")
			_, ok := e.Data.(event.ReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			return nil
		},
		Matcher:      domain.ReceiptCreatedEvent,
		AbortOnError: true,
	}
}
