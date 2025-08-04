package contact_event

import (
	"context"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	contact_repo "erp/project/core/contact/repository"
)

type contactEventHandler struct {
	contactRepo contact_repo.ContactRepository
	emitLog     logger.EmitLog
}

func NewContactEventHandler(
	logger logger.Logger,
	bus bus.Bus,
	contactRepo contact_repo.ContactRepository,
) {
	h := contactEventHandler{
		emitLog:     logger.EmitLog("contact-event"),
		contactRepo: contactRepo,
	}
	bus.RegisterHandler(domain.EventCustomerCreated, h.OnCreateCustomer())
	bus.RegisterHandler(domain.EventCustomerEdited, h.OnEditedCustomer())
	bus.RegisterHandler(domain.SupplierCreatedEvent, h.OnCreateSupplier())
	bus.RegisterHandler(domain.SupplierEditedEvent, h.OnEditedSupplier())
	bus.RegisterHandler(domain.DealEditedEvent, h.OnEditDeal())
	bus.RegisterHandler(domain.DealCreatedEvent, h.OnCreateDeal())
}

func (h *contactEventHandler) OnEditDeal() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnEditDeal"))
				}
			}()
			payload, ok := e.Data.(event.DealEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			payload.Req.SetContext(ctx)
			err = h.contactRepo.ContactBulkTx(payload.Tx, &payload.Req, dto.ContactBulkData{
				PartyID:  payload.Data.ID,
				Contacts: payload.Data.Contacts,
			})
			return
		},
		AbortOnError: true,
		Matcher:      domain.DealEditedEvent,
	}
}

func (h *contactEventHandler) OnCreateDeal() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnCreateDeal"))
				}
			}()
			payload, ok := e.Data.(event.DealEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			payload.Req.SetContext(ctx)
			err = h.contactRepo.ContactBulkTx(payload.Tx, &payload.Req, dto.ContactBulkData{
				PartyID:  payload.Deal.ID,
				Contacts: payload.Data.Contacts,
			})
			return
		},
		AbortOnError: true,
		Matcher:      domain.DealCreatedEvent,
	}
}

func (h *contactEventHandler) OnCreateSupplier() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnCreateSupplier"))
				}
			}()
			payload, ok := e.Data.(event.SupplierEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			payload.Req.SetContext(ctx)
			err = h.contactRepo.ContactBulkTx(payload.Tx, &payload.Req, dto.ContactBulkData{
				PartyID:  payload.Supplier.ID,
				Contacts: payload.Data.Contacts,
			})
			return nil
		},
		AbortOnError: true,
		Matcher:      domain.SupplierCreatedEvent,
	}
}

func (h *contactEventHandler) OnEditedSupplier() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnEditedSupplier"))
				}
			}()
			payload, ok := e.Data.(event.SupplierEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			payload.Req.SetContext(ctx)
			err = h.contactRepo.ContactBulkTx(payload.Tx, &payload.Req, dto.ContactBulkData{
				PartyID:  payload.Data.ID,
				Contacts: payload.Data.Contacts,
			})
			return nil
		},
		AbortOnError: true,
		Matcher:      domain.SupplierEditedEvent,
	}
}

func (h *contactEventHandler) OnCreateCustomer() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnCreateCustomer"))
				}
			}()
			payload, ok := e.Data.(event.CustomerEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			payload.Req.SetContext(ctx)
			err = h.contactRepo.ContactBulkTx(payload.Tx, &payload.Req, dto.ContactBulkData{
				PartyID:  payload.Customer.ID,
				Contacts: payload.Data.Contacts,
			})
			return nil
		},
		AbortOnError: true,
		Matcher:      domain.EventCustomerCreated,
	}
}

func (h *contactEventHandler) OnEditedCustomer() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnEditedCustomer"))
				}
			}()
			payload, ok := e.Data.(event.CustomerEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			payload.Req.SetContext(ctx)
			err = h.contactRepo.ContactBulkTx(payload.Tx, &payload.Req, dto.ContactBulkData{
				PartyID:  payload.Data.ID,
				Contacts: payload.Data.Contacts,
			})
			return nil
		},
		AbortOnError: true,
		Matcher:      domain.EventCustomerEdited,
	}
}
