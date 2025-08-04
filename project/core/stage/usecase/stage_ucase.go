package stage_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/logger"
	stage_repo "erp/project/core/stage/repository"
)

type StageUseCase interface {
	CreateStage(req *common.RequestContext, d dto.StageData) (res dto.StageDto, err error)
	EditStage(req *common.RequestContext, d dto.StageData) (err error)
	GetStages(req *common.RequestContext, d dto.StagesRequest) ( res dto.ResponseDataList[[]dto.StageDto],err error )
	StageTransition(req *common.RequestContext,d dto.StageTransitionData)(err error)
	DeleteStage(req *common.RequestContext,d *dto.DeleteRequest)(err error)
}

type stageUseCase struct {
	permission repository.PermissionService
	core       repository.CoreService
	repo       stage_repo.StageRepository
	emitLog    logger.EmitLog
}

func NewStageUseCase(
	permission repository.PermissionService,
	core repository.CoreService,
	repo stage_repo.StageRepository,
	logger logger.Logger,
	bus bus.Bus,
) StageUseCase {
	return &stageUseCase{
		permission: permission,
		core:       core,
		repo:       repo,
		emitLog:    logger.EmitLog("stage-usecase"),
	}
}

func(u *stageUseCase)DeleteStage(req *common.RequestContext,d *dto.DeleteRequest)(err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("DeleteStage"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.STAGE, domain.DELETE); err != nil {
		return err
	}
	err = u.repo.DeleteStage(req,d)
	return
}

func (u *stageUseCase) StageTransition(req *common.RequestContext,d dto.StageTransitionData)(err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("StageTransition"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.STAGE, domain.EDIT); err != nil {
		return err
	}
	err = u.repo.StageTransition(req, d)
	if err != nil {
		return
	}
	return
}

func (u *stageUseCase) CreateStage(req *common.RequestContext, d dto.StageData) (res dto.StageDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateStage"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.STAGE, domain.CREATE); err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}
	stage, err := u.repo.CreateStage(req, d)
	if err != nil {
		return
	}
	res = dto.StageFromModel(stage)
	return
}
func (u *stageUseCase) EditStage(req *common.RequestContext, d dto.StageData) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditStage"))
		}
	}()

	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.STAGE, domain.EDIT); err != nil {
		return err
	}
	err = u.repo.EditStage(req, d)
	if err != nil {
		return
	}

	return
}


func (u *stageUseCase) GetStages(req *common.RequestContext, d dto.StagesRequest) (
	res dto.ResponseDataList[[]dto.StageDto],err error ){
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetStages"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.STAGE, domain.VIEW); err != nil {
		return res, err
	}
	res.Body.Result, err = u.repo.GetStages(req, d)
	if err != nil {
		return
	}
	return
}

