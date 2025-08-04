package deal_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/logger"
	deal_repo "erp/project/crm/deal/repository"
)

type DealUseCase interface {
	CreateDeal(req *common.RequestContext, d dto.DealData) (res dto.DealDto, err error)
	EditDeal(req *common.RequestContext, d dto.DealData) (err error)
	GetDeal(req *common.RequestContext, d dto.RequestEntity) (res dto.ResultEntity[dto.DealDetailDto], err error)
	GetDeals(req *common.RequestContext, d dto.DealsRequest) (res dto.ResponseDataList[[]dto.DealDto],
		err error)
	DealTransition(req *common.RequestContext, d dto.EntityTransitionData) (err error)
}

type dealUseCase struct {
	permission repository.PermissionService
	core       repository.CoreService
	repo       deal_repo.DealRepository
	emitLog    logger.EmitLog
	bus        bus.Bus
	c          di.Container
}

func NewDealUseCase(
	permission repository.PermissionService,
	core repository.CoreService,
	repo deal_repo.DealRepository,
	logger logger.Logger,
	bus bus.Bus,
	c di.Container,
) DealUseCase {
	return &dealUseCase{
		permission: permission,
		core:       core,
		repo:       repo,
		emitLog:    logger.EmitLog("deal-usecase"),
		bus:        bus,
		c:          c,
	}
}

func (u *dealUseCase) DealTransition(req *common.RequestContext, d dto.EntityTransitionData) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("DealTransition"))
		}
		err = u.closeTx(tx, err)
	}(tx)

	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.DEAL, domain.EDIT); err != nil {
		return err
	}
	err = u.repo.DealTransition(tx, req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.PartyStageChange, event.ChangeStageEventData{
		Tx:              tx,
		ProfileID:       req.Profile.ID,
		StageTransition: d,
	})
	return
}

func (u *dealUseCase) CreateDeal(req *common.RequestContext, d dto.DealData) (res dto.DealDto, err error) {
	tx := u.c.Get(domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateDeal"))
		}
		err = u.closeTx(tx, err)
	}(tx)

	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.DEAL, domain.CREATE); err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}
	deal, err := u.repo.CreateDeal(tx, req, d)
	if err != nil {
		return
	}
	res = dto.DealFromModel(deal)
	err = u.bus.Emit(req.Ctx, domain.DealCreatedEvent, event.DealEventData{
		Tx:   tx,
		Data: d,
		Deal: deal,
		Req:  *req,
	})
	return
}
func (u *dealUseCase) EditDeal(req *common.RequestContext, d dto.DealData) (err error) {
	
	tx := u.c.Get(domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditDeal"))
		}
		err = u.closeTx(tx, err)
	}(tx)

	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.DEAL, domain.EDIT); err != nil {
		return err
	}
	err = u.repo.EditDeal(tx, req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.DealEditedEvent, event.DealEventData{
		Tx:   tx,
		Data: d,
		Req:  *req,
	})
	return
}

func (u *dealUseCase) GetDeal(req *common.RequestContext, d dto.RequestEntity) (res dto.ResultEntity[dto.DealDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetDeal"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.DEAL, domain.VIEW); err != nil {
		return res, err
	}

	res.Entity, err = u.repo.GetDeal(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.Deal.ID)
	res.Contacts = u.core.GetPartyContacts(req, res.Entity.Deal.ID)
	res.Contacts = u.core.GetPartyContacts(req, res.Entity.Deal.ID)
	return
}

func (u *dealUseCase) GetDeals(req *common.RequestContext, d dto.DealsRequest) (res dto.ResponseDataList[[]dto.DealDto],
	err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetDeals"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.DEAL, domain.VIEW); err != nil {
		return res, err
	}

	res.Body.Result, err = u.repo.GetDeals(req, d)
	if err != nil {
		return
	}
	// res.Body.FilterOptions = u.repo.GetFilterOptions()
	return
}

func (u *dealUseCase) closeTx(tx *query.QueryTx, err error) error {
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
