package currency_exchange_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	currency_exchange_repo "erp/project/core/currency_exchange/internal/repository"
)

type CurrencyExchangeUcase interface {
	GetCurrencyExchange(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.CurrencyExchangeDto], err error)
	CreateCurrencyExchange(req *common.RequestContext, d *dto.CreateCurrencyExchangeRequest) (
		res dto.CurrencyExchangeDto, err error)
	GetCurrencyExchanges(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.CurrencyExchangeDto], err error)
	EditCurrencyExchange(req *common.RequestContext, d *dto.EditCurrencyExchangeRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
}

type currencyExchangeUcase struct {
	emitLog              logger.EmitLog
	currencyExchangeRepo currency_exchange_repo.CurrencyExchangeRepo
	permission           repository.PermissionService
	core                 repository.CoreService
	fsm                  fsm.FsmState
}

func NewCurrencyExchangeUcase(
	logger logger.Logger,
	currencyExchangeRepo currency_exchange_repo.CurrencyExchangeRepo,
	permission repository.PermissionService,
	core repository.CoreService,
	fsm fsm.FsmState,
) CurrencyExchangeUcase {
	return &currencyExchangeUcase{
		emitLog:              logger.EmitLog("currency-exchange-usecase"),
		currencyExchangeRepo: currencyExchangeRepo,
		permission:           permission,
		core:                 core,
		fsm:                  fsm,
	}
}

func (u *currencyExchangeUcase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CURRENCY_EXCHANGE, domain.EDIT); err != nil {
		return
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.currencyExchangeRepo.UpdateStatus(req, d, nextState)
	return
}

func (u *currencyExchangeUcase) EditCurrencyExchange(req *common.RequestContext, d *dto.EditCurrencyExchangeRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditCurrencyExchange"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CURRENCY_EXCHANGE, domain.EDIT)
	if err != nil {
		return
	}
	err = u.currencyExchangeRepo.EditCurrencyExchange(req, d)
	return
}

func (u *currencyExchangeUcase) GetCurrencyExchange(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.CurrencyExchangeDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCurrencyExchange"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CURRENCY_EXCHANGE, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.currencyExchangeRepo.GetCurrencyExchange(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
func (u *currencyExchangeUcase) CreateCurrencyExchange(req *common.RequestContext, d *dto.CreateCurrencyExchangeRequest) (
	res dto.CurrencyExchangeDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateCurrencyExchange"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CURRENCY_EXCHANGE, domain.CREATE)
	if err != nil {
		return
	}
	currencyExchange, err := u.currencyExchangeRepo.CreateCurrencyExchange(req, d)
	if err != nil {
		return
	}
	res = dto.CurrencyExchangeDtoFromModel(&currencyExchange)
	return
}
func (u *currencyExchangeUcase) GetCurrencyExchanges(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CurrencyExchangeDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCurrencyExchanges"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CURRENCY_EXCHANGE, domain.VIEW)
	if err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res, err = u.currencyExchangeRepo.GetCurrencyExchanges(req, d)
	if err != nil {
		return
	}
	return
}
