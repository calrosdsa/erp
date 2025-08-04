package cost_center_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	cost_center_repo "erp/project/accounting/cost_center/internal/repository"
)

type CostCenterUseCase interface {
	GetCostCenter(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.CostCenterDto], err error)
	CreateCostCenter(req *common.RequestContext, d *dto.CreateCostCenterRequet) (
		res dto.CostCenterDto, err error)
	GetCostCenters(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.CostCenterDto], err error)
}

type costCenterUcase struct {
	emitLog        logger.EmitLog
	costCenterRepo cost_center_repo.CostCenterRepository
	permission     repository.PermissionService
	core           repository.CoreService
}

func NewCostCenterUcase(
	logger logger.Logger,
	costCenterRepo cost_center_repo.CostCenterRepository,
	permission repository.PermissionService,
	core repository.CoreService,
) CostCenterUseCase {
	return &costCenterUcase{
		emitLog: logger.EmitLog("cost-center-usecase"),
		costCenterRepo: costCenterRepo,
		permission: permission,
		core: core,
	}
}

func (u *costCenterUcase) GetCostCenter(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.CostCenterDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCostCenter"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.COST_CENTER, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.costCenterRepo.GetCostCenter(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}
func (u *costCenterUcase) CreateCostCenter(req *common.RequestContext, d *dto.CreateCostCenterRequet) (
	res dto.CostCenterDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateCostCenter"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.COST_CENTER, domain.CREATE)
	if err != nil {
		return
	}
	res, err = u.costCenterRepo.CreateCostCenter(req, d)
	if err != nil {
		return
	}
	return
}
func (u *costCenterUcase) GetCostCenters(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CostCenterDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCostCenters"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.COST_CENTER, domain.CREATE)
	if err != nil {
		return
	}
	res, err = u.costCenterRepo.GetCostCenters(req, d)
	if err != nil {
		return
	}
	return
}
