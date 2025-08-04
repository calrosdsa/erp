package quotation_ucase

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
	quotation_repo "erp/project/quotation/repository"
	"fmt"
)

type QuotationUseCase interface {
	CreateQuotation(req *common.RequestContext, i dto.QuotationBody) (dto.QuotationDto, error)
	CreateQuotationTx(tx *query.QueryTx, req *common.RequestContext, i dto.QuotationBody) (res dto.QuotationDto, err error)
	GetQuotations(req *common.RequestContext, i *dto.RequestQuotations) (dto.PaginationResult[[]dto.QuotationDto], error)
	GetQuotation(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		dto.ResultEntity[dto.QuotationDetailDto], error,
	)
	EditQuotation(req *common.RequestContext, d dto.QuotationBody) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error

	GetQuotationEntity(partyCode string) (domain.EntityTemplate, error)
}

type quotationUsecase struct {
	quotationRepo     quotation_repo.QuotationRepository
	emitLog           logger.EmitLog
	permissionService repository.PermissionService
	invoiceFsm        fsm.FsmState
	bus               bus.Bus
	c                 di.Container
	core              repository.CoreService
	document repository.DocumentService
}

func NewQuotationUseCase(
	logger logger.Logger,
	permissionService repository.PermissionService,
	quotationRepo quotation_repo.QuotationRepository,
	invoiceFsm fsm.FsmState,
	bus bus.Bus,
	c di.Container,
	core repository.CoreService,
	document repository.DocumentService,
) QuotationUseCase {
	quotationUsecase := quotationUsecase{
		quotationRepo:     quotationRepo,
		emitLog:           logger.EmitLog("quotation-usecase"),
		permissionService: permissionService,
		invoiceFsm:        invoiceFsm,
		c:                 c,
		bus:               bus,
		core:              core,
		document: document,
	}
	return &quotationUsecase
}

func (u *quotationUsecase) UpdateStatus(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}()
	quotationEntity, err := u.GetQuotationEntity(i.Body.PartyType)
	if err != nil {
		return err
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, quotationEntity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := u.invoiceFsm.NextState(i.Body.CurrentState, i.Body.Events)
	if err != nil {
		return err
	}
	quotation, err := u.quotationRepo.UpdateStatus(tx, req, i.Body.PartyID, i.Body.CurrentState, nextState)
	if err != nil {
		return
	}
	lineItems, err := u.document.GetLineItems(tx.Query, req, quotation.ID)
	if err != nil {
		return
	}
	evtData := event.StatusQuotationEventData{
		Quotation:          quotation,
		Tx:                 tx,
		LineItemsData:      lineItems,
		QuotationPartyType: i.Body.PartyType,
	}
	switch nextState {
	case proto.State_UNPAID.String():
		err = u.bus.Emit(req.Ctx, domain.InvoiceSubmittedEvent, evtData)
	case proto.State_CANCELLED.String():
		err = u.bus.Emit(req.Ctx, domain.InvoiceCancelledEvent, evtData)
	}
	return
}

func (u *quotationUsecase) GetQuotation(req *common.RequestContext, i *dto.RequestEntityWithParty) (
	res dto.ResultEntity[dto.QuotationDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetQuotation"))
		}
	}()
	entity, err := u.GetQuotationEntity(i.PartyType)
	if err != nil {
		return
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED

	}
	res, err = u.quotationRepo.GetQuotation(req, i)
	if err != nil {
		return res, err
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.Quotation.ID)
	return
}

func (u *quotationUsecase) GetQuotations(req *common.RequestContext, i *dto.RequestQuotations) (
	res dto.PaginationResult[[]dto.QuotationDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetQuotations"))
		}
	}()
	entity, err := u.GetQuotationEntity(i.PartyType)
	if err != nil {
		return
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.quotationRepo.GetQuotations(req, i)
	if err != nil {
		return
	}
	res.FilterOptions = u.quotationRepo.GetFilterOptions()
	return
}

func (u *quotationUsecase) EditQuotation(req *common.RequestContext, d dto.QuotationBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditQuotation"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	entity, err := u.GetQuotationEntity(d.Quotation.QuotationPartyType)
	if err != nil {
		return err
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	if err = u.quotationRepo.EditQuotation(tx, req, d); err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.QuotationEditEvent, event.QuotationEventData{
		Tx:   tx,
		Body: d,
	})
	return
}

func (u *quotationUsecase) CreateQuotation(req *common.RequestContext, i dto.QuotationBody) (res dto.QuotationDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateQuotation"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	res, err = u.createQuotation(tx, req, i)
	return
}
func (u *quotationUsecase) CreateQuotationTx(tx *query.QueryTx, req *common.RequestContext,
	i dto.QuotationBody) (res dto.QuotationDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateQuotation"))
		}
	}()
	res, err = u.createQuotation(tx, req, i)
	return
}

func (u *quotationUsecase) createQuotation(tx *query.QueryTx, req *common.RequestContext,
	i dto.QuotationBody) (res dto.QuotationDto, err error) {
	entity, err := u.GetQuotationEntity(i.Quotation.QuotationPartyType)
	if err != nil {
		return res, err
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	quotation, err := u.quotationRepo.CreateQuotation(tx, req, i)
	if err != nil {
		return
	}
	res = dto.QuotationDtoFromModel(&quotation)
	err = u.bus.Emit(req.Ctx, domain.QuotationCreatedEvent, event.QuotationEventData{
		Tx:        tx,
		Body:      i,
		Quotation: &quotation,
	})
	return
}

func (u *quotationUsecase) closeTx(tx *query.QueryTx, err error) error {
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

func (u *quotationUsecase) GetQuotationEntity(partyCode string) (domain.EntityTemplate, error) {
	switch partyCode {
	case proto.PartyType_supplierQuotation.String():
		return domain.SUPPLIER_QUOTATION, nil
	case proto.PartyType_salesQuotation.String():
		return domain.QUOTATION, nil
	default:
		return domain.EntityTemplate{}, domain.PARTY_TYPE_NOT_FOUND
	}
}
