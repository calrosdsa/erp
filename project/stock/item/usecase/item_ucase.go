package item_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	item_repo "erp/project/stock/item/repository"
	itemprice_ucase "erp/project/stock/itemprice/usecase"
	"fmt"
	"strconv"
)

type ItemUseCase interface {
	GetItemDetail(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ItemDetailDto], err error)
	GetItems(req *common.RequestContext, i *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.ItemDto], err error)
	CreateItem(req *common.RequestContext, d dto.ItemData) (dto.ItemDto, error)
	CreateItemTx(query *query.QueryTx, req *common.RequestContext, d dto.ItemData) (dto.ItemDto, error)
	UpdateItem(req *common.RequestContext, i *dto.UpdateItemRequest) error
	EditItem(req *common.RequestContext, d dto.ItemData) error
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error)
}

type itemUseCase struct {
	emitLog    logger.EmitLog
	itemRepo   item_repo.ItemRepository
	permission repository.PermissionService
	core       repository.CoreService
	c          di.Container
	bus        bus.Bus
	fsm        fsm.FsmState

	itemPriceUcase itemprice_ucase.ItemPriceUseCase
}

func NewItemUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	core repository.CoreService,
	itemRepo item_repo.ItemRepository,
	c di.Container,
	bus bus.Bus,
	fsm fsm.FsmState,
) ItemUseCase {
	return &itemUseCase{
		emitLog:    logger.EmitLog("item-ucase"),
		permission: permission,
		itemRepo:   itemRepo,
		core:       core,
		c:          c,
		bus:        bus,
		fsm:        fsm,
		itemPriceUcase: c.Get(domain.ItemPriceUseCase).(itemprice_ucase.ItemPriceUseCase),
	}
}

func (u *itemUseCase) UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (
	err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.ITEM, domain.EDIT); err != nil {
		return err
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.itemRepo.UpdateStatus(req, d, nextState)
	return
}

func (u *itemUseCase) EditItem(req *common.RequestContext, d dto.ItemData) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditItem"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)

	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.ITEM, domain.VIEW)
	if err != nil {
		return err
	}
	err = u.itemRepo.EditItem(tx, req, d)
	if err != nil {
		return
	}

	err = u.bus.Emit(req.Ctx, domain.ItemEditedEvent, event.ItemCreatedEventData{
		Tx:   tx,
		Body: d,
		Req:  *req,
	})
	return
}

func (u *itemUseCase) GetItemDetail(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ItemDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemDetail"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.ITEM, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.itemRepo.GetItemDetail(req, d)
	if err != nil {
		return
	}
	itemPrices,err := u.itemPriceUcase.GetItemPrices(req,dto.ItemPricesRequest{
		ItemID: strconv.Itoa(int(res.Entity.ID)),
	})
	res.Entity.ItemPrices = itemPrices.Body.Result
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
func (u *itemUseCase) GetItems(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.ItemDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItems"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.ITEM, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.itemRepo.GetItems(req, i)
	if err != nil {
		return
	}
	return
}

func (u *itemUseCase) CreateItemTx(tx *query.QueryTx, req *common.RequestContext, d dto.ItemData) (
	res dto.ItemDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateItemTx"))
		}
	}()
	res, err = u.createItem(tx, req, d)
	return
}
func (u *itemUseCase) CreateItem(req *common.RequestContext, d dto.ItemData) (
	res dto.ItemDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateItem"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	res, err = u.createItem(tx, req, d)
	return
}

func (u *itemUseCase) createItem(tx *query.QueryTx, req *common.RequestContext, d dto.ItemData) (
	res dto.ItemDto, err error) {
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.ITEM, domain.CREATE); err != nil {
		return
	}
	res, err = u.itemRepo.CreateItem(tx, req, d)
	if err != nil {
		return
	}

	err = u.bus.Emit(req.Ctx, domain.ItemCreatedEvent, event.ItemCreatedEventData{
		Item: &res,
		Tx:   tx,
		Body: d,
		Req:  *req,
	})
	return
}

func (u *itemUseCase) UpdateItem(req *common.RequestContext, i *dto.UpdateItemRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateItem"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.ITEM, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.itemRepo.UpdateItem(req, i)
	return
}

func (s *itemUseCase) closeTx(tx *query.QueryTx, err error) error {
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
