package stock_entry_ucase

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
	stock_entry_repo "erp/project/stock/stock_entry/internal/repository"
	"fmt"
)

type StockEntryUseCase interface {
	GetStockEntry(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.StockEntryDto], err error)
	CreateStockEntry(req *common.RequestContext, d dto.StockEntryBody) (
		res dto.StockEntryDto, err error)
	GetStockEntries(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.StockEntryDto], err error)
	UpdateStockEntryStatus(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error)
	EditStockEntry(req *common.RequestContext, d dto.StockEntryBody) (err error)
}

type stockEntryUcase struct {
	emitLog        logger.EmitLog
	stockEntryRepo stock_entry_repo.StockEntryRepository
	permission     repository.PermissionService
	core           repository.CoreService
	fsm            fsm.FsmState
	c              di.Container
	bus            bus.Bus
	document repository.DocumentService
}

func NewStockEntrytUcase(
	logger logger.Logger,
	stockEntryRepo stock_entry_repo.StockEntryRepository,
	permission repository.PermissionService,
	core repository.CoreService,
	c di.Container,
	bus bus.Bus,
	fsm fsm.FsmState,
	document repository.DocumentService,
) StockEntryUseCase {
	return &stockEntryUcase{
		emitLog:        logger.EmitLog("stock-entry-usecase"),
		stockEntryRepo: stockEntryRepo,
		permission:     permission,
		fsm:            fsm,
		core:           core,
		c:              c,
		bus:            bus,
		document:          document,
	}
}

func (u *stockEntryUcase) EditStockEntry(req *common.RequestContext, d dto.StockEntryBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditStockEntry"))
		}
		err = u.closeTx(tx, err)
	}(tx)

	if allow := u.permission.CheckPermission(req.Ctx, req, domain.STOCK_ENTRY, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	if err = u.stockEntryRepo.EditStockEntry(tx, req, d); err != nil {
		return
	}
	u.bus.Emit(req.Ctx, domain.StockEntryEditEvent, event.StockEntryEventData{
		Tx:             tx,
		StockEntryBody: d,
	})

	return
}

func (u *stockEntryUcase) UpdateStockEntryStatus(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpstateOrderState"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.STOCK_ENTRY, domain.EDIT)
	nextState, err := u.fsm.NextState(i.Body.CurrentState, i.Body.Events, i.Body.PartyType)
	if err != nil {
		return err
	}
	stockEntry, err := u.stockEntryRepo.UpdateStockEntryStatus(req, tx, i.Body.PartyID, i.Body.CurrentState, nextState)
	if err != nil {
		return
	}
	lineItems, err := u.document.GetLineItems(tx.Query, req, stockEntry.ID,
		repository.OptionsStock.WithLoadItemInLine(true),
		repository.OptionsStock.WithLoadLineStockEntry(true),
	)
	if err != nil {
		return
	}
	fmt.Println("LINE ITEMS", lineItems)
	payload := event.StatusStockEntryEventData{
		Tx:            tx,
		StockEntry:    stockEntry,
		CompanyID:     req.ActiveCompany.ID,
		LineItemsData: lineItems,
	}
	switch nextState {
	case proto.State_SUBMITTED.String():
		err = u.bus.Emit(req.Ctx, domain.StockEntrySubmittedEvent, payload)
	}
	return
}

func (u *stockEntryUcase) GetStockEntry(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.StockEntryDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetStockEntry"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.STOCK_ENTRY, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.stockEntryRepo.GetStockEntry(req, d)
	if err!= nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}
func (u *stockEntryUcase) CreateStockEntry(req *common.RequestContext, d dto.StockEntryBody) (
	res dto.StockEntryDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateStockEntry"))
		}
		err1 := u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err1)
	}(tx)

	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.STOCK_ENTRY, domain.CREATE)
	if err != nil {
		return
	}
	stockEntry, err := u.stockEntryRepo.CreateStockEntry(tx, req, d)
	if err != nil {
		return
	}
	payload := event.StockEntryEventData{
		StockEntry:     stockEntry,
		Tx:             tx,
		StockEntryBody: d,
	}
	res = dto.StockEntryDtoFromModel(&stockEntry)
	err = u.bus.Emit(req.Ctx, domain.EventStockEntryCreated, payload)
	return
}
func (u *stockEntryUcase) GetStockEntries(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.StockEntryDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetStockEntries"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.STOCK_ENTRY, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.stockEntryRepo.GetStockEntries(req, d)
	if err != nil {
		return
	}
	return
}

func (s *stockEntryUcase) closeTx(tx *query.QueryTx, err error) error {
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
