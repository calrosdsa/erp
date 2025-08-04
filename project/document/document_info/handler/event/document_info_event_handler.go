package documentinfo_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	documentinfo_repo "erp/project/document/document_info/repository"
)

type documentInfoEvent struct {
	emitLog        logger.EmitLog
	dInfoEventRepo documentinfo_repo.DocumentInfoEventRepo
}

func NewDocumentEventHandler(
	logger        logger.Logger,
	dInfoEventRepo documentinfo_repo.DocumentInfoEventRepo,
	bus bus.Bus,
) {
	h := documentInfoEvent{
		emitLog: logger.EmitLog("document-info-event"),
		dInfoEventRepo: dInfoEventRepo,
	}
	bus.RegisterHandler(domain.OrderCreatedEvent,h.OnCreatedOrder())
	bus.RegisterHandler(domain.ReceiptCreatedEvent,h.OnCreatedReceipt())
	bus.RegisterHandler(domain.InvoiceCreatedEvent,h.OnCreatedInvoice())
	bus.RegisterHandler(domain.CashOutflowCreatedEvent,h.OnCreateCashOutflow())
	// bus.RegisterHandler(domain.InvoiceCreatedEvent,h.OnCash())
}
func (h *documentInfoEvent) OnCreateCashOutflow() bus.Handler {
	return  bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func ()  {
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnCreateCashOutflow"))
				}
			}()
			payload,ok:= e.Data.(event.CashOutflowEventData)
			if !ok  {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.dInfoEventRepo.CreateDocumentInfoForCashOutflow(ctx,payload)
			return 
		},
		AbortOnError: true,
		Matcher: domain.CashOutflowCreatedEvent,
	}
}

func (h *documentInfoEvent) OnCreatedReceipt() bus.Handler {
	return  bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func ()  {
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnCreatedReceipt"))
				}
			}()
			payload,ok:= e.Data.(event.ReceiptEventData)
			if !ok  {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.dInfoEventRepo.CreateDocumentInfoForReceipt(ctx,payload)
			return 
		},
		AbortOnError: true,
		Matcher: domain.ReceiptCreatedEvent,
	}
}

func (h *documentInfoEvent) OnCreatedInvoice() bus.Handler {
	return  bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func ()  {
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnCreatedInvoice"))
				}
			}()
			payload,ok:= e.Data.(event.InvoiceEventData)
			if !ok  {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.dInfoEventRepo.CreateDocumentInfoForInvoice(ctx,payload)
			return 
		},
		AbortOnError: true,
		Matcher: domain.InvoiceCreatedEvent,
	}
}

func (h *documentInfoEvent) OnCreatedOrder() bus.Handler {
	return  bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func ()  {
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnCreatedOrder"))
				}
			}()
			payload,ok:= e.Data.(event.OrderEventData)
			if !ok  {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.dInfoEventRepo.CreateDocumentInfoForOrder(ctx,payload)
			return 
		},
		AbortOnError: true,
		Matcher: domain.OrderCreatedEvent,
	}
}


