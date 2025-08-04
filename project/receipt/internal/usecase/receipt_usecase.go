package receipt_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain"
	"fmt"

	// "erp/internal/domain/event"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	receipt_repo "erp/project/receipt/internal/repository"
	// "fmt"
)

type ReceiptUseCase interface {
	CreateReceipt(req *common.RequestContext, i dto.ReceiptBody) (dto.ReceiptDto, error)
	GetReceipts(req *common.RequestContext, i *dto.RequestReceipts) (
		dto.PaginationResult[[]dto.ReceiptDto], error)
	GetReceiptDetail(req *common.RequestContext, i *dto.RequestEntityWithParty) (
		dto.ResultEntity[dto.ReceiptDetailDto], error)

	GetReceiptEntity(partyType string) (domain.EntityTemplate, error)
	EditReceipt(req *common.RequestContext, d dto.ReceiptBody) (err error)
	UpdateReceiptState(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
	ExportReceipt(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error)
}

type receiptUseCase struct {
	receiptRepo receipt_repo.ReceiptRepository
	permission  repository.PermissionService
	emitLog     logger.EmitLog
	bus         bus.Bus
	c           di.Container
	core        repository.CoreService
	stock       repository.StockService
	document repository.DocumentService
}

func NewReceiptUseCase(
	permission repository.PermissionService,
	logger logger.Logger,
	receiptRepo receipt_repo.ReceiptRepository,
	bus bus.Bus,
	c di.Container,
	core repository.CoreService,
	stock repository.StockService,
	document repository.DocumentService,
) ReceiptUseCase {
	return &receiptUseCase{
		permission:  permission,
		emitLog:     logger.EmitLog("receipt-usecase"),
		receiptRepo: receiptRepo,
		bus:         bus,
		c:           c,
		core:        core,
		stock:       stock,
		document: document,
	}
}

func (u *receiptUseCase) ExportReceipt(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ExportReceipt"))
		}
	}()
	order,err := u.GetReceiptDetail(req,&dto.RequestEntityWithParty{
		ID: i.ID,
		PartyType: i.PartyType,
	})
	if err != nil {
		return
	}
	res,err  =u.document.GenerateReceiptDocumentPdf(req,order.Entity.Receipt,i.PartyType)
	return
}

func (u *receiptUseCase) GetReceipts(req *common.RequestContext, i *dto.RequestReceipts) (
	res dto.PaginationResult[[]dto.ReceiptDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetReceipts"))
		}

	}()
	entity, err := u.GetReceiptEntity(i.PartyType)
	if err != nil {
		return
	}

	if allow := u.permission.CheckPermission(req.Ctx, req, entity, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.receiptRepo.GetReceipts(req, i)
	if err != nil {
		return
	}
	res.FilterOptions = u.receiptRepo.GetFilterOptions()
	return
}

func (u *receiptUseCase) GetReceiptDetail(req *common.RequestContext, i *dto.RequestEntityWithParty) (
	res dto.ResultEntity[dto.ReceiptDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetReceiptDetail"))
		}
	}()
	entity, err := u.GetReceiptEntity(i.PartyType)
	if err != nil {
		return
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entity, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.receiptRepo.GetReceiptDetail(req, i)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.Receipt.ID)
	return
}

func (u *receiptUseCase) UpdateReceiptState(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateReceiptState"))
		}
		err1 := u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err1)
	}(tx)
	receiptEntity, err := u.GetReceiptEntity(i.Body.PartyType)
	if err != nil {
		return domain.PARTY_TYPE_NOT_FOUND
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, receiptEntity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	// fmt.Println("RECEIPT STATE",i.Body.CurrenctState,i.Body.Events)
	nextState, err := u.nextState(i.Body.CurrentState, i.Body.Events)
	if err != nil {
		return err
	}
	receipt, err := u.receiptRepo.UpdateReceiptState(ctx, req, i.Body.PartyID, i.Body.CurrentState, nextState)
	if err != nil {
		return
	}
	var opt repository.OptionStock
	if proto.PartyType_purchaseReceipt.String() == i.Body.PartyType {
		opt = repository.OptionsStock.WithLoadReceiptLineItem(true)
	}
	if proto.PartyType_deliveryNote.String() == i.Body.PartyType {
		opt = repository.OptionsStock.WithLoadDeliveryLineItem(true)
	}
	lineItemsData, err := u.document.GetLineItems(tx.Query, req, receipt.ID,
		repository.OptionsStock.WithLoadItemInLine(true),
		opt,
	)
	if err != nil {
		return
	}
	//Emit event for subscribers
	itemDefault, err := u.stock.GetStockDefault(req)
	payload := event.StatusReceiptEventData{
		Receipt:          *receipt,
		LineItemsData:    lineItemsData,
		Tx:               tx,
		Company:          req.ActiveCompany,
		ReceiptPartyType: i.Body.PartyType,
		CompanyDefault:   req.CompanyDefaults,
		StockDefault:     itemDefault,
	}
	switch nextState {
	case proto.State_TO_BILL.String():
		err = u.bus.Emit(req.Ctx, domain.ReceiptSubmittedEvent, payload)
	case proto.State_CANCELLED.String():
		err = u.bus.Emit(req.Ctx, domain.ReceiptCancelledEvent, payload)
	}
	return
}

func (u *receiptUseCase) EditReceipt(req *common.RequestContext, d dto.ReceiptBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditReceipt"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	entity, err := u.GetReceiptEntity(d.Receipt.ReceiptPartyType)
	if err != nil {
		return err
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	if err = u.receiptRepo.EditReceipt(tx, req, d); err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.ReceiptEditEvent, event.ReceiptEventData{
		Tx:   tx,
		Body: d,
	})
	return
}

func (u *receiptUseCase) CreateReceipt(req *common.RequestContext, i dto.ReceiptBody) (
	res dto.ReceiptDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateReceipt"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	entity, err := u.GetReceiptEntity(i.Receipt.ReceiptPartyType)
	if err != nil {
		return res, domain.PARTY_TYPE_NOT_FOUND
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entity, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	receipt, err := u.receiptRepo.CreateReceipt(req, tx, i)
	if err != nil {
		return
	}
	res = dto.ReceiptDtoFromModel(&receipt)
	//Emit event for subscribers
	fmt.Println("SENDING EVENT RECEIPT ...")
	u.bus.Emit(req.Ctx, domain.ReceiptCreatedEvent, event.ReceiptEventData{
		Tx:               tx,
		Body:             i,
		ReceiptPartyType: i.Receipt.ReceiptPartyType,
		Receipt:          receipt,
	})

	// u.bus.Emit(req.Ctx, domain.EventInvoiceCreated, event.CreateInvoiceEventData{
	// 	Tx: tx,
	// 	CreateInvoiceBody:i,
	// 	Invoice:invoice,
	// })
	return
}

func (u *receiptUseCase) nextState(state string, events []int32) (string, error) {
	eventsState := make([]proto.EventState, len(events))
	for i, event := range events {
		eventsState[i] = proto.EventState(event)
	}
	stateMachine := fsm.New()
	draftState := stateMachine.NewState(proto.State_DRAFT.String())
	toBill := stateMachine.NewState(proto.State_TO_BILL.String())
	cancellledState := stateMachine.NewState(proto.State_CANCELLED.String())
	paidState := stateMachine.NewState(proto.State_PAID.String())

	submit := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_SUBMIT_EVENT))
	cancel := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_CANCEL_EVENT))

	stateMachine.LinkStates(draftState, toBill, submit)
	stateMachine.LinkStates(toBill, cancellledState, cancel)
	stateMachine.LinkStates(paidState, cancellledState, cancel)

	err := stateMachine.SetInialState(state)
	if err != nil {
		return "", err
	}

	stateMachine.Compute(eventsState, true)
	if currentState, ok := stateMachine.PresentState.Value.(string); ok {
		return currentState, nil
	} else {
		return "", domain.FAIL_TYPE_ASSERTION
	}
}

func (u *receiptUseCase) closeTx(tx *query.QueryTx, err error) error {
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

func (u *receiptUseCase) GetReceiptEntity(partyType string) (domain.EntityTemplate, error) {
	switch partyType {
	case proto.PartyType_purchaseReceipt.String():
		return domain.PURCHASE_RECEIPT, nil
	case proto.PartyType_deliveryNote.String():
		return domain.DELIVERY_NOTE, nil
	default:
		return domain.EntityTemplate{}, domain.PARTY_NOT_FOUND
	}
}
