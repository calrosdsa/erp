package module_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	module_repo "erp/project/core/module/repository"
)

type ModuleUsecase interface {
	GetModule(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ModuleDetailDto], err error)
	GetModuleDetail(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ModuleDetailDto], err error)
	CreateModule(req *common.RequestContext, d dto.ModuleData) (
		res dto.ModuleDto, err error)
	GetModules(req *common.RequestContext, d dto.ModulesRequest) (
		res []dto.ModuleDto, err error)
	EditModule(req *common.RequestContext, d dto.ModuleData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
	GetEntitiesSearch(req *common.RequestContext, d *dto.ModuleSearchRequest) (
		res []dto.EntityDto, err error,
	)
}

type moduleUcase struct {
	emitLog    logger.EmitLog
	moduleRepo module_repo.ModuleRepository
	permission repository.PermissionService
	core       repository.CoreService
	fsm        fsm.FsmState
}

func NewModuleUcase(
	logger logger.Logger,
	moduleRepo module_repo.ModuleRepository,
	permission repository.PermissionService,
	core repository.CoreService,
	fsm fsm.FsmState,
) ModuleUsecase {
	return &moduleUcase{
		emitLog:    logger.EmitLog("module-usecase"),
		moduleRepo: moduleRepo,
		permission: permission,
		core:       core,
		fsm:        fsm,
	}
}
func (u *moduleUcase) GetEntitiesSearch(req *common.RequestContext, d *dto.ModuleSearchRequest) (
	res []dto.EntityDto, err error,
) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetEntitiesSearch"))
		}
	}()
	res, err = u.moduleRepo.GetEntitiesSearch(req, d)
	return
}

func (u *moduleUcase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.MODULE, domain.EDIT); err != nil {
		return
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.moduleRepo.UpdateStatus(req, d, nextState)
	return
}

func (u *moduleUcase) EditModule(req *common.RequestContext, d dto.ModuleData) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditModule"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.MODULE, domain.EDIT)
	if err != nil {
		return
	}
	err = u.moduleRepo.EditModule(req, d)
	return
}

func (u *moduleUcase) GetModuleDetail(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ModuleDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetModuleDetail"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.MODULE, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.moduleRepo.GetModuleDetail(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.Module.ID)
	return
}

func (u *moduleUcase) GetModule(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ModuleDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetModule"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.MODULE, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.moduleRepo.GetModule(req, d)
	if err != nil {
		return
	}
	return
}
func (u *moduleUcase) CreateModule(req *common.RequestContext, d dto.ModuleData) (
	res dto.ModuleDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateModule"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.MODULE, domain.CREATE)
	if err != nil {
		return
	}
	module, err := u.moduleRepo.CreateModule(req, d)
	if err != nil {
		return
	}
	res = dto.ModuleDtoFromModel(module)
	return
}
func (u *moduleUcase) GetModules(req *common.RequestContext, d dto.ModulesRequest) (
	res []dto.ModuleDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetModules"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.MODULE, domain.VIEW)
	if err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res, err = u.moduleRepo.GetModules(req, d)
	if err != nil {
		return
	}
	return
}
