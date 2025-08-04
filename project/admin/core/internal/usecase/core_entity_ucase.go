package core_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	pg_core "erp/project/admin/core/internal/repository"
)

type CoreEntityUseCase interface {
	GetEntities(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.EntityDto], error)
	GetEntity(req *common.AdminRequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.EntityDetailDto], error)
	CreateEntity(req *common.AdminRequestContext, d *dto.CreateEntityRequest) error
	AddEntityAction(req *common.AdminRequestContext, d *dto.AddEntityActionRequest) error
}

type coreEntityUseCase struct {
	emitLog        logger.EmitLog
	coreEntityRepo pg_core.CoreEntityRepository
}

func NewCoreEntityUseCase(
	logger logger.Logger,
	coreEntityRepo pg_core.CoreEntityRepository,
) CoreEntityUseCase {
	return &coreEntityUseCase{
		emitLog:        logger.EmitLog("core-entity-usecase"),
		coreEntityRepo: coreEntityRepo,
	}
}
func (u *coreEntityUseCase) AddEntityAction(req *common.AdminRequestContext, d *dto.AddEntityActionRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("AddEntityAction"))
		}
	}()
	err = u.coreEntityRepo.AddEntityAction(req, d)
	return
}

func (u *coreEntityUseCase) GetEntity(req *common.AdminRequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.EntityDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetEntity"))
		}
	}()
	res, err = u.coreEntityRepo.GetEntity(req, d)
	return
}

func (u *coreEntityUseCase) CreateEntity(req *common.AdminRequestContext, d *dto.CreateEntityRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateEntity"))
		}
	}()
	err = u.coreEntityRepo.CreateEntity(req, d)
	return
}

func (u *coreEntityUseCase) GetEntities(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.EntityDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetEntities"))
		}
	}()
	res, err = u.coreEntityRepo.GetEntities(req, d)
	return
}
