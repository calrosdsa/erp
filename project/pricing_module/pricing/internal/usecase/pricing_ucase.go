package pricing_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	pricing_repo "erp/project/pricing_module/pricing/internal/repository"
)

type PricingUseCase interface {
	GetPricing(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.PricingDetailDto], err error)
	CreatePricing(req *common.RequestContext, d *dto.CreatePricingRequest) (
		res dto.PricingDto, err error)
	GetPricings(req *common.RequestContext, d *dto.RequestPricings) (
		res dto.PaginationResult[[]dto.PricingDto], err error)
	EditPricing(req *common.RequestContext, d *dto.EditPricingRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
}

type pricingUcase struct {
	emitLog     logger.EmitLog
	pricingRepo pricing_repo.PricingRepository
	permission  repository.PermissionService
	core        repository.CoreService
	fsm         fsm.FsmState
}

func NewPricingUcase(
	logger logger.Logger,
	pricingRepo pricing_repo.PricingRepository,
	permission repository.PermissionService,
	core repository.CoreService,
	fsm fsm.FsmState,
) PricingUseCase {
	return &pricingUcase{
		emitLog:     logger.EmitLog("currency-exchange-usecase"),
		pricingRepo: pricingRepo,
		permission:  permission,
		core:        core,
		fsm:         fsm,
	}
}

func (u *pricingUcase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICING, domain.EDIT); err != nil {
		return
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.pricingRepo.UpdateStatus(req, d, nextState)
	return
}

func (u *pricingUcase) EditPricing(req *common.RequestContext, d *dto.EditPricingRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditPricing"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICING, domain.EDIT)
	if err != nil {
		return
	}
	err = u.pricingRepo.EditPricing(req, d)
	return
}

func (u *pricingUcase) GetPricing(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.PricingDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPricing"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICING, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.pricingRepo.GetPricing(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.PricingDto.ID)
	return
}
func (u *pricingUcase) CreatePricing(req *common.RequestContext, d *dto.CreatePricingRequest) (
	res dto.PricingDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePricing"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICING, domain.CREATE)
	if err != nil {
		return
	}
	pricing, err := u.pricingRepo.CreatePricing(req, d)
	if err != nil {
		return
	}
	res = dto.PricingDtoFromModel(&pricing)
	return
}
func (u *pricingUcase) GetPricings(req *common.RequestContext, d *dto.RequestPricings) (
	res dto.PaginationResult[[]dto.PricingDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPricings"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICING, domain.VIEW)
	if err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res, err = u.pricingRepo.GetPricings(req, d)
	if err != nil {
		return
	}
	return
}
