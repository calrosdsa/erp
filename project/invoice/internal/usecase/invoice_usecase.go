package invoice_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	invoice_repo "erp/project/invoice/internal/repository"
	"fmt"
)

type InvoiceUseCase interface {
	CreateInvoice(req *common.RequestContext, i dto.InvoiceBody) (dto.InvoiceDto, error)
	GetInvoices(req *common.RequestContext, i *dto.RequestInvoices) (dto.PaginationResult[[]dto.InvoiceDto], error)
	GetInvoiceDetail(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		dto.ResultEntity[dto.InvoiceDetailDto], error,
	)
	UpdateInvoiceState(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
	EditInvoice(req *common.RequestContext, d dto.InvoiceBody) (err error)
	GetEntityInvoice(partyCode string) (domain.EntityTemplate, error)

	ExportInvoice(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error)
}

type invoiceUseCase struct {
	invoiceRepository invoice_repo.InvoiceRepository
	emitLog           logger.EmitLog
	permissionService repository.PermissionService
	invoiceFsm        fsm.FsmState
	bus               bus.Bus
	c                 di.Container
	core              repository.CoreService
	document          repository.DocumentService
}

func NewInvoiceUseCase(
	logger logger.Logger,
	permissionService repository.PermissionService,
	invoiceRepository invoice_repo.InvoiceRepository,
	invoiceFsm fsm.FsmState,
	bus bus.Bus,
	c di.Container,
	core repository.CoreService,
	document repository.DocumentService,
) InvoiceUseCase {
	invoiceUseCase := invoiceUseCase{
		invoiceRepository: invoiceRepository,
		emitLog:           logger.EmitLog("invoice-usecase"),
		permissionService: permissionService,
		invoiceFsm:        invoiceFsm,
		c:                 c,
		bus:               bus,
		core:              core,
		document:          document,
	}
	return &invoiceUseCase
}

func (u *invoiceUseCase) ExportInvoice(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ExportInvoice"))
		}
	}()
	order, err := u.GetInvoiceDetail(req, &dto.RequestEntityWithParty{
		ID:        i.ID,
		PartyType: i.PartyType,
	})
	if err != nil {
		return
	}
	res, err = u.document.GenerateInvoiceDocumentPdf(req, order.Entity.Invoice, i.PartyType)
	return
}

func (u *invoiceUseCase) UpdateInvoiceState(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateInvoiceState"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}()
	invoiceEntity, err := u.GetEntityInvoice(i.Body.PartyType)
	if err != nil {
		return err
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, invoiceEntity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := u.invoiceFsm.NextState(i.Body.CurrentState, i.Body.Events)
	if err != nil {
		return err
	}
	invoice, err := u.invoiceRepository.UpdateInvoiceState(tx, req, i.Body.PartyID, i.Body.CurrentState, nextState)
	if err != nil {
		return
	}
	var opt repository.OptionStock
	fmt.Println("INVOICE LINE OPTION", i.Body.PartyType, invoice.UpdateStock)
	if proto.PartyType_purchaseInvoice.String() == i.Body.PartyType && invoice.UpdateStock {
		opt = repository.OptionsStock.WithLoadReceiptLineItem(true)
	}
	if proto.PartyType_saleInvoice.String() == i.Body.PartyType && invoice.UpdateStock {
		opt = repository.OptionsStock.WithLoadDeliveryLineItem(true)
	}
	lineItemsData, err := u.document.GetLineItems(tx.Query, req, invoice.ID,
		repository.OptionsStock.WithLoadItemInLine(true),
		opt,
	)
	if err != nil {
		return
	}
	taxLinesData, err := u.document.GetTaxLines(tx.Query, req.Ctx, invoice.ID)
	if err != nil {
		return
	}
	evtData := event.StatusInvoiceEventData{
		Invoice:          *invoice,
		Tx:               tx,
		LineItemsData:    lineItemsData,
		TaxLinesData:     taxLinesData,
		InvoicePartyType: i.Body.PartyType,
		CompanyDefaults:  req.CompanyDefaults,
	}
	switch nextState {
	case proto.State_UNPAID.String():
		err = u.bus.Emit(req.Ctx, domain.InvoiceSubmittedEvent, evtData)
	case proto.State_CANCELLED.String():
		err = u.bus.Emit(req.Ctx, domain.InvoiceCancelledEvent, evtData)
	}
	return
}

func (u *invoiceUseCase) EditInvoice(req *common.RequestContext, d dto.InvoiceBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditInvoice"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	entity, err := u.GetEntityInvoice(d.Invoice.InvoicePartyType)
	if err != nil {
		return err
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	if err = u.invoiceRepository.EditInvoice(tx, req, d); err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.InvoiceEditEvent, event.InvoiceEventData{
		Tx:   tx,
		Body: d,
	})
	return
}

func (u *invoiceUseCase) GetInvoiceDetail(req *common.RequestContext, i *dto.RequestEntityWithParty) (
	res dto.ResultEntity[dto.InvoiceDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetInvoiceDetail"))
		}
	}()
	entity, err := u.GetEntityInvoice(i.PartyType)
	if err != nil {
		return
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.invoiceRepository.GetInvoiceDetail(req, i)
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.Invoice.ID)
	return
}

func (u *invoiceUseCase) GetInvoices(req *common.RequestContext, i *dto.RequestInvoices) (
	res dto.PaginationResult[[]dto.InvoiceDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetInvoices"))
		}
	}()
	entity, err := u.GetEntityInvoice(i.PartyType)
	if err != nil {
		return
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res, err = u.invoiceRepository.GetInvoices(req, i)
	if err != nil {
		return
	}
	res.FilterOptions = u.invoiceRepository.GetFilterOptions(string(req.LanguageCode), i.PartyType)
	return
}

func (u *invoiceUseCase) CreateInvoice(req *common.RequestContext, i dto.InvoiceBody) (res dto.InvoiceDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateInvoice"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	fmt.Println(i.Invoice.InvoicePartyType)
	entity, err := u.GetEntityInvoice(i.Invoice.InvoicePartyType)
	if err != nil {
		return res, err
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	invoice, err := u.invoiceRepository.CreateInvoice(req, tx, i)
	if err != nil {
		return
	}
	res = dto.InvoiceDtoFromModel(&invoice)
	err = u.bus.Emit(req.Ctx, domain.InvoiceCreatedEvent, event.InvoiceEventData{
		Tx:               tx,
		Body:             i,
		InvoicePartyType: i.Invoice.InvoicePartyType,
		Invoice:          invoice,
	})
	return
}

func (u *invoiceUseCase) closeTx(tx *query.QueryTx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}

func (u *invoiceUseCase) GetEntityInvoice(partyCode string) (domain.EntityTemplate, error) {
	switch partyCode {
	case proto.PartyType_purchaseInvoice.String():
		return domain.PURCHASE_INVOICE, nil
	case proto.PartyType_saleInvoice.String():
		return domain.SALE_INVOICE, nil
	default:
		return domain.EntityTemplate{}, domain.PARTY_TYPE_NOT_FOUND
	}
}
