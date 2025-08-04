package itemline_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	itemline_repo "erp/project/stock/itemline/internal/repository"
	"fmt"
)

type ItemLineEventHandler struct {
	itemLineRepo itemline_repo.ItemLineEventRepo
	emitLog      logger.EmitLog
}

func NewItemLineEventHandler(
	bus bus.Bus,
	logger logger.Logger,
	itemLineRepo itemline_repo.ItemLineEventRepo,
) {
	handler := ItemLineEventHandler{
		itemLineRepo: itemLineRepo,
		emitLog:      logger.EmitLog("item-line-events"),
	}
	fmt.Println("REGISTER ITEM LINE EVENTS...")
	bus.RegisterHandler(domain.InvoiceCreatedEvent, handler.OnInvoiceCreated())
	bus.RegisterHandler(domain.OrderCreatedEvent, handler.OnOrderCreated())
	bus.RegisterHandler(domain.ReceiptCreatedEvent, handler.OnReceiptCreated())
	bus.RegisterHandler(domain.EventStockEntryCreated, handler.OnStockEntryCreated())
	bus.RegisterHandler(domain.QuotationCreatedEvent, handler.OnQuotationCreated())

	bus.RegisterHandler(domain.InvoiceEditEvent, handler.OnInvoiceEdited())
	bus.RegisterHandler(domain.OrderEditEvent, handler.OnOrderEdited())
	bus.RegisterHandler(domain.ReceiptEditEvent, handler.OnReceiptEdited())
	bus.RegisterHandler(domain.StockEntryEditEvent, handler.OnStockEntryEdited())
	bus.RegisterHandler(domain.QuotationEditEvent, handler.OnQuotationEdited())

}

func (h *ItemLineEventHandler) OnQuotationEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnQuotationEdited"))
				}
			}()
			payload, ok := e.Data.(event.QuotationEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.EditQuotationLineItems(ctx, payload)
			return 
		},
		Matcher:      domain.QuotationEditEvent,
		AbortOnError: true,
	}
}

func (h *ItemLineEventHandler) OnStockEntryEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnStockEntryEdited"))
				}
			}()
			payload, ok := e.Data.(event.StockEntryEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.EditStockEntryLineItems(ctx, payload)
			return 
		},
		Matcher:      domain.StockEntryEditEvent,
		AbortOnError: true,
	}
}

func (h *ItemLineEventHandler) OnReceiptEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnReceiptEdited"))
				}
			}()
			payload, ok := e.Data.(event.ReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.EditReceiptLineItems(ctx, payload)
			return 
		},
		Matcher:      domain.ReceiptEditEvent,
		AbortOnError: true,
	}
}

func (h *ItemLineEventHandler) OnOrderEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnOrderEdited"))
				}
			}()
			payload, ok := e.Data.(event.OrderEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.EditOrderLineItems(ctx, payload)
			return 
		},
		Matcher:      domain.OrderEditEvent,
		AbortOnError: true,
	}
}

func (h *ItemLineEventHandler) OnInvoiceEdited() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnInvoiceEdited"))
				}
			}()
			payload, ok := e.Data.(event.InvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.EditInvoiceLineItems(ctx, payload)
			return 
		},
		Matcher:      domain.InvoiceEditEvent,
		AbortOnError: true,
	}
}


func (h *ItemLineEventHandler) OnQuotationCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnQuotationCreated"))
				}
			}()
			payload, ok := e.Data.(event.QuotationEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.CreateQuotationLineItems(ctx, payload)
			return 
		},
		Matcher:      domain.QuotationCreatedEvent,
		AbortOnError: true,
	}
}

func (h *ItemLineEventHandler) OnStockEntryCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnStockEntryCreated"))
				}
			}()
			payload, ok := e.Data.(event.StockEntryEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.CreateStockEntryItemLines(ctx, payload)
			return 
		},
		Matcher:      domain.EventStockEntryCreated,
		AbortOnError: true,
	}
}

func (h *ItemLineEventHandler) OnReceiptCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnReceiptCreated"))
				}
			}()
			payload, ok := e.Data.(event.ReceiptEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.CreateReceiptItemLines(ctx, payload)
			return 
		},
		Matcher:      domain.ReceiptCreatedEvent,
		AbortOnError: true,
	}
}

func (h *ItemLineEventHandler) OnOrderCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnOrderCreated"))
				}
			}()
			payload, ok := e.Data.(event.OrderEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.CreateOrderItemLines(ctx, payload)
			return 
		},
		Matcher:      domain.OrderCreatedEvent,
		AbortOnError: true,
	}
}

func (h *ItemLineEventHandler) OnInvoiceCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnInvoiceCreated"))
				}
			}()
			payload, ok := e.Data.(event.InvoiceEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.itemLineRepo.CreateInvoiceItemLines(ctx, payload)
			return 
		},
		Matcher:      domain.InvoiceCreatedEvent,
		AbortOnError: true,
	}
}
